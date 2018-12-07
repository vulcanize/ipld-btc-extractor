// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vow_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/drip_file/vow"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/shared_behaviors"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Drip file vow repository", func() {
	var (
		db                    *postgres.DB
		dripFileVowRepository vow.DripFileVowRepository
		headerRepository      datastore.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
		dripFileVowRepository = vow.DripFileVowRepository{}
		dripFileVowRepository.SetDB(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.DripFileVowModel
		modelWithDifferentLogIdx.LogIndex++
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.DripFileVowChecked,
			LogEventTableName:        "maker.drip_file_vow",
			TestModel:                test_data.DripFileVowModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &dripFileVowRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("adds a drip file vow event", func() {
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
			err = dripFileVowRepository.Create(headerID, []interface{}{test_data.DripFileVowModel})

			Expect(err).NotTo(HaveOccurred())
			var dbDripFileVow vow.DripFileVowModel
			err = db.Get(&dbDripFileVow, `SELECT what, data, log_idx, tx_idx, raw_log FROM maker.drip_file_vow WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbDripFileVow.What).To(Equal(test_data.DripFileVowModel.What))
			Expect(dbDripFileVow.Data).To(Equal(test_data.DripFileVowModel.Data))
			Expect(dbDripFileVow.LogIndex).To(Equal(test_data.DripFileVowModel.LogIndex))
			Expect(dbDripFileVow.TransactionIndex).To(Equal(test_data.DripFileVowModel.TransactionIndex))
			Expect(dbDripFileVow.Raw).To(MatchJSON(test_data.DripFileVowModel.Raw))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MissingHeadersBehaviorInputs{
			Repository:    &dripFileVowRepository,
			RepositoryTwo: &vow.DripFileVowRepository{},
		}

		shared_behaviors.SharedRepositoryMissingHeadersBehaviors(&inputs)
	})

	Describe("MissingHeaders", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.DripFileVowChecked,
			Repository:              &dripFileVowRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})