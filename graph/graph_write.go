package graph

import (
	"bufio"
	"github.com/KevinSmall/ethgraph/logr"
	"github.com/yaricom/goGraphML/graphml"
	"os"
)

func WriteGraph(filename string, gr *graphml.GraphML) error {
	file, err := os.Create(filename)
	if err != nil {
		logr.Error.Panicln(err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	err = gr.Encode(writer, false)
	return err
}
