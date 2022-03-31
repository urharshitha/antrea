// Copyright 2020 Antrea Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package addressgroup

import (
	"io"
	"reflect"
	"sort"
	"time"

	"antrea.io/antrea/pkg/antctl/transform"
	"antrea.io/antrea/pkg/antctl/transform/common"
	cpv1beta "antrea.io/antrea/pkg/apis/controlplane/v1beta2"
)

type Response struct {
	Name string               `json:"name" yaml:"name"`
	Pods []common.GroupMember `json:"pods,omitempty"`
}

func listTransform(l interface{}, opts map[string]string) (interface{}, error) {
	groups := l.(*cpv1beta.AddressGroupList)
	sortBy := ""
	if sb, ok := opts["sort-by"]; ok {
		sortBy = sb
	}
	adsorter := &Adsorter{
		addressgroups: groups.Items,
		sortBy:        sortBy,
	}
	sort.Sort(adsorter)

	result := make([]Response, 0, len(groups.Items))
	for i := range adsorter.addressgroups {
		o, _ := objectTransform(&adsorter.addressgroups[i], opts)
		result = append(result, o.(Response))
	}
	return result, nil
}

func objectTransform(o interface{}, _ map[string]string) (interface{}, error) {
	group := o.(*cpv1beta.AddressGroup)
	var pods []common.GroupMember
	for _, pod := range group.GroupMembers {
		pods = append(pods, common.GroupMemberPodTransform(pod))
	}
	return Response{Name: group.Name, Pods: pods}, nil
}

func Transform(reader io.Reader, single bool, opts map[string]string) (interface{}, error) {
	return transform.GenericFactory(
		reflect.TypeOf(cpv1beta.AddressGroup{}),
		reflect.TypeOf(cpv1beta.AddressGroupList{}),
		objectTransform,
		listTransform,
		opts,
	)(reader, single)
}

const sortBycreationtime = "CreationTimestamp"

type TimeSlice []time.Time
type Adsorter struct {
	addressgroups []cpv1beta.AddressGroup
	sortBy        string
}

func (ads *Adsorter) Len() int { return len(ads.addressgroups) }
func (ads *Adsorter) Swap(i, j int) {
	ads.addressgroups[i].CreationTimestamp, ads.addressgroups[j].CreationTimestamp = ads.addressgroups[j].CreationTimestamp, ads.addressgroups[i].CreationTimestamp
}

func (ads *Adsorter) Less(i, j int) bool {
	switch ads.sortBy {
	case sortBycreationtime:
		return ads.addressgroups[i].CreationTimestamp.Before(&ads.addressgroups[j].CreationTimestamp)
	default:
		return ads.addressgroups[i].Name < ads.addressgroups[j].Name
	}
}

var _ common.TableOutput = new(Response)

func (r Response) GetTableHeader() []string {
	return []string{"NAME", "POD-IPS"}
}

func (r Response) GetPodNames(maxColumnLength int) string {
	list := make([]string, len(r.Pods))
	for i, pod := range r.Pods {
		list[i] = pod.IP
	}
	return common.GenerateTableElementWithSummary(list, maxColumnLength)
}

func (r Response) GetTableRow(maxColumnLength int) []string {
	return []string{r.Name, r.GetPodNames(maxColumnLength)}
}

func (r Response) SortRows() bool {
	return true
}
