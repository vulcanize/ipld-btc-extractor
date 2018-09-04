package vat_init_test

import (
	"database/sql"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_init"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("", func() {
	Describe("Create", func() {
		It("adds a vat event", func() {
			db := test_config.NewTestDB(core.Node{})
			test_config.CleanTestDB(db)
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, err := headerRepository.CreateOrUpdateHeader(core.Header{})
			Expect(err).NotTo(HaveOccurred())
			vatInitRepository := vat_init.NewVatInitRepository(db)

			err = vatInitRepository.Create(headerID, test_data.VatInitModel)

			Expect(err).NotTo(HaveOccurred())
			var dbVatInit vat_init.VatInitModel
			err = db.Get(&dbVatInit, `SELECT ilk,tx_idx, raw_log FROM maker.vat_init WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbVatInit.Ilk).To(Equal(test_data.VatInitModel.Ilk))
			Expect(dbVatInit.TransactionIndex).To(Equal(test_data.VatInitModel.TransactionIndex))
			Expect(dbVatInit.Raw).To(MatchJSON(test_data.VatInitModel.Raw))
		})

		It("does not duplicate vat events", func() {
			db := test_config.NewTestDB(core.Node{})
			test_config.CleanTestDB(db)
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, err := headerRepository.CreateOrUpdateHeader(core.Header{})
			Expect(err).NotTo(HaveOccurred())
			vatInitRepository := vat_init.NewVatInitRepository(db)
			err = vatInitRepository.Create(headerID, test_data.VatInitModel)
			Expect(err).NotTo(HaveOccurred())

			err = vatInitRepository.Create(headerID, test_data.VatInitModel)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("pq: duplicate key value violates unique constraint"))
		})

		It("removes vat if corresponding header is deleted", func() {
			db := test_config.NewTestDB(core.Node{})
			test_config.CleanTestDB(db)
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, err := headerRepository.CreateOrUpdateHeader(core.Header{})
			Expect(err).NotTo(HaveOccurred())
			vatInitRepository := vat_init.NewVatInitRepository(db)
			err = vatInitRepository.Create(headerID, test_data.VatInitModel)
			Expect(err).NotTo(HaveOccurred())

			_, err = db.Exec(`DELETE FROM headers WHERE id = $1`, headerID)

			Expect(err).NotTo(HaveOccurred())
			var dbVatInit vat_init.VatInitModel
			err = db.Get(&dbVatInit, `SELECT ilk, tx_idx, raw_log FROM maker.vat_init WHERE header_id = $1`, headerID)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(sql.ErrNoRows))
		})
	})

	Describe("MissingHeaders", func() {
		It("returns headers with no associated vat event", func() {
			db := test_config.NewTestDB(core.Node{})
			test_config.CleanTestDB(db)
			headerRepository := repositories.NewHeaderRepository(db)
			startingBlockNumber := int64(1)
			vatInitBlockNumber := int64(2)
			endingBlockNumber := int64(3)
			blockNumbers := []int64{startingBlockNumber, vatInitBlockNumber, endingBlockNumber, endingBlockNumber + 1}
			var headerIDs []int64
			for _, n := range blockNumbers {
				headerID, err := headerRepository.CreateOrUpdateHeader(core.Header{BlockNumber: n})
				headerIDs = append(headerIDs, headerID)
				Expect(err).NotTo(HaveOccurred())
			}
			vatInitRepository := vat_init.NewVatInitRepository(db)
			err := vatInitRepository.Create(headerIDs[1], test_data.VatInitModel)
			Expect(err).NotTo(HaveOccurred())

			headers, err := vatInitRepository.MissingHeaders(startingBlockNumber, endingBlockNumber)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(2))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber)))
		})

		It("only returns headers associated with the current node", func() {
			db := test_config.NewTestDB(core.Node{})
			test_config.CleanTestDB(db)
			blockNumbers := []int64{1, 2, 3}
			headerRepository := repositories.NewHeaderRepository(db)
			dbTwo := test_config.NewTestDB(core.Node{ID: "second"})
			headerRepositoryTwo := repositories.NewHeaderRepository(dbTwo)
			var headerIDs []int64
			for _, n := range blockNumbers {
				headerID, err := headerRepository.CreateOrUpdateHeader(core.Header{BlockNumber: n})
				Expect(err).NotTo(HaveOccurred())
				headerIDs = append(headerIDs, headerID)
				_, err = headerRepositoryTwo.CreateOrUpdateHeader(core.Header{BlockNumber: n})
				Expect(err).NotTo(HaveOccurred())
			}
			vatInitRepository := vat_init.NewVatInitRepository(db)
			vatInitRepositoryTwo := vat_init.NewVatInitRepository(dbTwo)
			err := vatInitRepository.Create(headerIDs[0], test_data.VatInitModel)
			Expect(err).NotTo(HaveOccurred())

			nodeOneMissingHeaders, err := vatInitRepository.MissingHeaders(blockNumbers[0], blockNumbers[len(blockNumbers)-1])
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeOneMissingHeaders)).To(Equal(len(blockNumbers) - 1))

			nodeTwoMissingHeaders, err := vatInitRepositoryTwo.MissingHeaders(blockNumbers[0], blockNumbers[len(blockNumbers)-1])
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeTwoMissingHeaders)).To(Equal(len(blockNumbers)))
		})
	})
})