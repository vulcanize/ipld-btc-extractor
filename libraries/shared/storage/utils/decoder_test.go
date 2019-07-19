// VulcanizeDB
// Copyright © 2019 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package utils_test

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
)

var _ = Describe("Storage decoder", func() {
	It("decodes uint256", func() {
		fakeInt := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000539")
		row := utils.StorageDiffRow{StorageValue: fakeInt}
		metadata := utils.StorageValueMetadata{Type: utils.Uint256}

		result, err := utils.Decode(row, metadata)

		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(big.NewInt(0).SetBytes(fakeInt.Bytes()).String()))
	})

	It("decodes uint128", func() {
		fakeInt := common.HexToHash("0000000000000000000000000000000000000000000000000000000000011123")
		row := utils.StorageDiffRow{StorageValue: fakeInt}
		metadata := utils.StorageValueMetadata{Type: utils.Uint128}

		result, err := utils.Decode(row, metadata)

		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(big.NewInt(0).SetBytes(fakeInt.Bytes()).String()))
	})

	It("decodes uint48", func() {
		fakeInt := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000123")
		row := utils.StorageDiffRow{StorageValue: fakeInt}
		metadata := utils.StorageValueMetadata{Type: utils.Uint48}

		result, err := utils.Decode(row, metadata)

		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(big.NewInt(0).SetBytes(fakeInt.Bytes()).String()))
	})

	It("decodes address", func() {
		fakeAddress := common.HexToAddress("0x12345")
		row := utils.StorageDiffRow{StorageValue: fakeAddress.Hash()}
		metadata := utils.StorageValueMetadata{Type: utils.Address}

		result, err := utils.Decode(row, metadata)

		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(fakeAddress.Hex()))
	})

	Describe("when there are multiple items packed in the storage slot", func() {
		It("decodes uint48 items", func() {
			//this is a real storage data example
			packedStorage := common.HexToHash("000000000000000000000000000000000000000000000002a300000000002a30")
			row := utils.StorageDiffRow{StorageValue: packedStorage}
			packedTypes := map[int]utils.ValueType{}
			packedTypes[0] = utils.Uint48
			packedTypes[1] = utils.Uint48

			metadata := utils.StorageValueMetadata{
				Type:        utils.PackedSlot,
				PackedTypes: packedTypes,
			}

			result, err := utils.Decode(row, metadata)
			decodedValues := result.(map[int]string)

			Expect(err).NotTo(HaveOccurred())
			Expect(decodedValues[0]).To(Equal(big.NewInt(0).SetBytes(common.HexToHash("2a30").Bytes()).String()))
			Expect(decodedValues[1]).To(Equal(big.NewInt(0).SetBytes(common.HexToHash("2a300").Bytes()).String()))
		})

		It("decodes 5 uint48 items", func() {
			//TODO: this packedStorageHex was generated by hand, it would be nice to test this against
			//real storage data that has several items packed into it
			packedStorageHex := "0000000A5D1AFFFFFFFFFFFE00000009F3C600000002A300000000002A30"

			packedStorage := common.HexToHash(packedStorageHex)
			row := utils.StorageDiffRow{StorageValue: packedStorage}
			packedTypes := map[int]utils.ValueType{}
			packedTypes[0] = utils.Uint48
			packedTypes[1] = utils.Uint48
			packedTypes[2] = utils.Uint48
			packedTypes[3] = utils.Uint48
			packedTypes[4] = utils.Uint48

			metadata := utils.StorageValueMetadata{
				Type:        utils.PackedSlot,
				PackedTypes: packedTypes,
			}

			result, err := utils.Decode(row, metadata)
			decodedValues := result.(map[int]string)

			Expect(err).NotTo(HaveOccurred())
			Expect(decodedValues[0]).To(Equal(big.NewInt(0).SetBytes(common.HexToHash("2a30").Bytes()).String()))
			Expect(decodedValues[1]).To(Equal(big.NewInt(0).SetBytes(common.HexToHash("2a300").Bytes()).String()))
			Expect(decodedValues[2]).To(Equal(big.NewInt(0).SetBytes(common.HexToHash("9F3C6").Bytes()).String()))
			Expect(decodedValues[3]).To(Equal(big.NewInt(0).SetBytes(common.HexToHash("FFFFFFFFFFFE").Bytes()).String()))
			Expect(decodedValues[4]).To(Equal(big.NewInt(0).SetBytes(common.HexToHash("A5D1A").Bytes()).String()))
		})

		It("decodes 2 uint128 items", func() {
			//TODO: this packedStorageHex was generated by hand, it would be nice to test this against
			//real storage data that has several items packed into it
			packedStorageHex := "000000038D7EA4C67FF8E502B6730000" +
				"0000000000000000AB54A98CEB1F0AD2"
			packedStorage := common.HexToHash(packedStorageHex)
			row := utils.StorageDiffRow{StorageValue: packedStorage}
			packedTypes := map[int]utils.ValueType{}
			packedTypes[0] = utils.Uint128
			packedTypes[1] = utils.Uint128

			metadata := utils.StorageValueMetadata{
				Type:        utils.PackedSlot,
				PackedTypes: packedTypes,
			}

			result, err := utils.Decode(row, metadata)
			decodedValues := result.(map[int]string)

			Expect(err).NotTo(HaveOccurred())
			Expect(decodedValues[0]).To(Equal(big.NewInt(0).SetBytes(common.HexToHash("AB54A98CEB1F0AD2").Bytes()).String()))
			Expect(decodedValues[1]).To(Equal(big.NewInt(0).SetBytes(common.HexToHash("38D7EA4C67FF8E502B6730000").Bytes()).String()))
		})
	})
})
