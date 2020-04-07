// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
//
// Copyright 2019 Centrifuge GmbH
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

package types

var ExamplaryExtrinsic = Extrinsic{Version: 0x84, Signature: ExtrinsicSignatureV4{Signer: Address{IsAccountID: true, AsAccountID: AccountID{0xd4, 0x35, 0x93, 0xc7, 0x15, 0xfd, 0xd3, 0x1c, 0x61, 0x14, 0x1a, 0xbd, 0x4, 0xa9, 0x9f, 0xd6, 0x82, 0x2c, 0x85, 0x58, 0x85, 0x4c, 0xcd, 0xe3, 0x9a, 0x56, 0x84, 0xe7, 0xa5, 0x6d, 0xa2, 0x7d}, IsAccountIndex: false, AsAccountIndex: 0x0}, Signature: MultiSignature{IsSr25519: true, AsSr25519: Signature{0x5c, 0x77, 0x1d, 0xd5, 0x6a, 0xe0, 0xce, 0xed, 0x68, 0xd, 0xb3, 0xbb, 0x4c, 0x40, 0x7a, 0x38, 0x96, 0x99, 0x97, 0xae, 0xb6, 0xa, 0x2c, 0x62, 0x39, 0x1, 0x6, 0x2f, 0x7f, 0x8e, 0xbf, 0x2f, 0xe7, 0x73, 0x3a, 0x61, 0x3c, 0xf1, 0x6b, 0x78, 0xf6, 0x10, 0xc6, 0x52, 0x32, 0xa2, 0x3c, 0xc5, 0xce, 0x25, 0xda, 0x29, 0xa3, 0xd5, 0x84, 0x85, 0xd8, 0x7b, 0xd8, 0x3d, 0xb8, 0x18, 0x3f, 0x8}}, Era: ExtrinsicEra{IsImmortalEra: true, IsMortalEra: false, AsMortalEra: MortalEra{First: 0x0, Second: 0x0}}, Nonce: 0x1, Tip: 0x2}, Method: Call{CallIndex: CallIndex{SectionIndex: 0x3, MethodIndex: 0x0}, Args: Args{0xff, 0x8e, 0xaf, 0x4, 0x15, 0x16, 0x87, 0x73, 0x63, 0x26, 0xc9, 0xfe, 0xa1, 0x7e, 0x25, 0xfc, 0x52, 0x87, 0x61, 0x36, 0x93, 0xc9, 0x12, 0x90, 0x9c, 0xb2, 0x26, 0xaa, 0x47, 0x94, 0xf2, 0x6a, 0x48, 0xe5, 0x6c}}} //nolint:lll

var ExamplaryExtrinsicAccountID = ExtrinsicAccountID{Version: 0x84, Signature: ExtrinsicSignatureV4AccountId{Signer: AccountID{0xd4, 0x35, 0x93, 0xc7, 0x15, 0xfd, 0xd3, 0x1c, 0x61, 0x14, 0x1a, 0xbd, 0x4, 0xa9, 0x9f, 0xd6, 0x82, 0x2c, 0x85, 0x58, 0x85, 0x4c, 0xcd, 0xe3, 0x9a, 0x56, 0x84, 0xe7, 0xa5, 0x6d, 0xa2, 0x7d}, Signature: MultiSignature{IsSr25519: true, AsSr25519: Signature{0x5c, 0x77, 0x1d, 0xd5, 0x6a, 0xe0, 0xce, 0xed, 0x68, 0xd, 0xb3, 0xbb, 0x4c, 0x40, 0x7a, 0x38, 0x96, 0x99, 0x97, 0xae, 0xb6, 0xa, 0x2c, 0x62, 0x39, 0x1, 0x6, 0x2f, 0x7f, 0x8e, 0xbf, 0x2f, 0xe7, 0x73, 0x3a, 0x61, 0x3c, 0xf1, 0x6b, 0x78, 0xf6, 0x10, 0xc6, 0x52, 0x32, 0xa2, 0x3c, 0xc5, 0xce, 0x25, 0xda, 0x29, 0xa3, 0xd5, 0x84, 0x85, 0xd8, 0x7b, 0xd8, 0x3d, 0xb8, 0x18, 0x3f, 0x8}}, Era: ExtrinsicEra{IsImmortalEra: true, IsMortalEra: false, AsMortalEra: MortalEra{First: 0x0, Second: 0x0}}, Nonce: 0x1, Tip: 0x2}, Method: Call{CallIndex: CallIndex{SectionIndex: 0x3, MethodIndex: 0x0}, Args: Args{0xff, 0x8e, 0xaf, 0x4, 0x15, 0x16, 0x87, 0x73, 0x63, 0x26, 0xc9, 0xfe, 0xa1, 0x7e, 0x25, 0xfc, 0x52, 0x87, 0x61, 0x36, 0x93, 0xc9, 0x12, 0x90, 0x9c, 0xb2, 0x26, 0xaa, 0x47, 0x94, 0xf2, 0x6a, 0x48, 0xe5, 0x6c}}} //nolint:lll
