package trecresults

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

// The result file contains a map of all qrels lists, indexed by topic ID.
type QrelsFile struct {
	Qrels map[string]Qrels
}

// Qrels is a map of docids to relevance value
type Qrels map[string]*Qrel

type Qrel struct {
	Topic     string // The topic that this qrel is associated with
	Iteration string // Ignored by treceval
	DocId     string // the docid
	Score     int64  // the relevance score for this document
}

// Constructor for a QrelsFile pointer
func NewQrelsFile() *QrelsFile {
	return &QrelsFile{make(map[string]Qrels)}
}

// Creates a qrel structure from a single line from a results file.
//
// Returns parsing errors if any of the integer or float fields do not parse.
//
// Returns an error if there are not 4 fields in the result line.
//
// On error, a nil result is returned.
func QrelFromLine(line string) (*Qrel, error) {
	split := strings.Fields(line)

	if len(split) != 4 {
		err := errors.New("Incorrect number of fields in qrel string: " + line)
		return nil, err
	}

	topic := split[0]
	iteration := split[1]
	docId := split[2]

	score, err := strconv.ParseInt(split[3], 10, 0)
	if err != nil {
		return nil, err
	}
	return &Qrel{topic, iteration, docId, score}, nil
}

// This function returns a QrelsFile object created from the
// provided reader (eg a file).
//
// On errors, a QrelsFile containing all qrels read before the error was encountered is
// returned, along with the error.
func QrelsFromReader(file io.Reader) (QrelsFile, error) {
	var qf QrelsFile
	qf.Qrels = make(map[string]Qrels)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		q, err := QrelFromLine(scanner.Text())
		qrels, ok := qf.Qrels[q.Topic]
		if !ok {
			qrels = make(map[string]*Qrel)
			qf.Qrels[q.Topic] = qrels
		}
		if err != nil {
			return qf, err
		}
		qf.Qrels[q.Topic][q.DocId] = q
	}

	if err := scanner.Err(); err != nil {
		return qf, err
	}
	return qf, nil
}
