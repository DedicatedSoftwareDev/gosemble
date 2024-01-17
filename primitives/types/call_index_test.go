package types

// TODO:
//var call = Call{
//	CallIndex: CallIndex{
//		ModuleIndex:   0,
//		FunctionIndex: 0,
//	},
//	Args: []goscale.Encodable{},
//}
//
//func Test_NewCall(t *testing.T) {
//	var testExamples = []struct {
//		label       string
//		input       Call
//		expectation Call
//	}{
//		{
//			label: "Encode(Call(System.remark(0xab, 0xcd)))",
//			input: call,
//			expectation: Call{
//				CallIndex: CallIndex{
//					ModuleIndex:   0,
//					FunctionIndex: 0,
//				},
//				Args: []goscale.Encodable{},
//			},
//		},
//	}
//
//	for _, testExample := range testExamples {
//		t.Run(testExample.label, func(t *testing.T) {
//			assert.Equal(t, testExample.input.CallIndex.ModuleIndex, testExample.expectation.CallIndex.ModuleIndex)
//			assert.Equal(t, testExample.input.CallIndex.FunctionIndex, testExample.expectation.CallIndex.FunctionIndex)
//			assert.Equal(t, testExample.input.Args, testExample.expectation.Args)
//		})
//	}
//}
//
//func Test_EncodeCall(t *testing.T) {
//	var testExamples = []struct {
//		label       string
//		input       Call
//		expectation []byte
//	}{
//		{
//			label:       "Encode(Call(System.remark(0xab, 0xcd)))",
//			input:       call,
//			expectation: []byte{0x0, 0x0},
//		},
//	}
//
//	for _, testExample := range testExamples {
//		t.Run(testExample.label, func(t *testing.T) {
//			buffer := &bytes.Buffer{}
//
//			testExample.input.Encode(buffer)
//
//			assert.Equal(t, testExample.expectation, buffer.Bytes())
//		})
//	}
//}
//
//func Test_DecodeCall(t *testing.T) {
//	var testExamples = []struct {
//		label       string
//		input       []byte
//		expectation Call
//	}{
//		{
//			label:       "Decode(0x0, 0x0, 0x8, 0xab, 0xcd)",
//			input:       []byte{0x0, 0x0},
//			expectation: call,
//		},
//	}
//
//	for _, testExample := range testExamples {
//		t.Run(testExample.label, func(t *testing.T) {
//			buffer := &bytes.Buffer{}
//			buffer.Write(testExample.input)
//
//			result := DecodeCall(buffer)
//
//			assert.Equal(t, testExample.expectation, result)
//		})
//	}
//}
