package serializer

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
)

// WriteProtobufToBinaryFile write proto buffer message to binary file
func WriteProtobufToBinaryFile(message proto.Message, filename string) error {
	data, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("cannot marshal proto message to binary:%v", err)
	}

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("cannot write binary data to file:%w", err)
	}

	return nil
}

// ReadProtobufFromBinaryFile read protocol buffer from binary file
func ReadProtobufFromBinaryFile(filename string, message proto.Message) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("cannot read binary data from file:%w", err)
	}

	err = proto.Unmarshal(data, message)
	if err != nil {
		fmt.Errorf("cannot unmarshal binary to proto message:%w", err)
	}

	return nil
}
