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

import ( //"io"
	//"reflect"
	//"sort"
	//"antrea.io/antrea/pkg/antctl/transform"
	"antrea.io/antrea/pkg/antctl/transform/common"
	cpv1beta "antrea.io/antrea/pkg/apis/controlplane/v1beta2"
)

type Response struct {
	Name  string               `json:"name" yaml:"name"`
	Pods  []common.GroupMember `json:"pods,omitempty"`
	Nodes []common.GroupMember `json:"nodes,omitempty"`
}

type Adsorter struct {
	Addressgroups []cpv1beta.AddressGroup
	SortBy        string
}

var _ common.TableOutput = new(Response)

func (r Response) GetTableHeader() []string {
	return []string{"NAME", "POD-IPS", "NODE-IPS"}
}

func (r Response) GetPodIPs(maxColumnLength int) string {
	list := make([]string, len(r.Pods))
	for i, pod := range r.Pods {
		list[i] = pod.IP
	}
	return common.GenerateTableElementWithSummary(list, maxColumnLength)
}

func (r Response) GetNodeIPs(maxColumnLength int) string {
	list := make([]string, len(r.Nodes))
	for i, node := range r.Nodes {
		list[i] = node.IP
	}
	return common.GenerateTableElementWithSummary(list, maxColumnLength)
}

func (r Response) GetTableRow(maxColumnLength int) []string {
	return []string{r.Name, r.GetPodIPs(maxColumnLength), r.GetNodeIPs(maxColumnLength)}
}

func (r Response) SortRows() bool {
	return true
}
