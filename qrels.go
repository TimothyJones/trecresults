// Package trecresults provides helper functions for reading and writing trec results files
// suitable for using with treceval.
// 
// It has three main concepts:
//
// ResultFile: Contains a map of results for all topics contained in this results file
//
// ResultList: A slice containing the results for this topic
//
// Result: The data that describes a single entry in a result list
package trecresults

import (
  "errors"
  "strconv"
  "strings"
  "io"
  "bufio"
)

// The result file contains a map of all qrels lists, indexed by topic ID.
type QrelsFile struct {
  Qrels map[int64]Qrels
}

// Qrels is a map of docids to relevance value
type Qrels map[string]*Qrel

type Qrel struct {
  Topic int64 // The topic that this qrel is associated with
  Iteration string // Ignored by treceval
  DocId string // the docid
  Score int64 // the relevance score for this document
}

// Constructor for a QrelsFile pointer
func NewQrelsFile() *QrelsFile{
  return &QrelsFile{make(map[int64]Qrels)}
}

// Creates a result structure from a single line from a results file.
//
// Returns parsing errors if any of the integer or float fields do not parse.
//
// Returns an error if there are not 6 fields in the result line.
//
// On error, a nil result is returned.
// 201 0 AP880221-0047 0
// 201 0 AP880223-0069 0
func QrelFromLine(line string) (*Qrel, error) {
  split := strings.Fields(line)

  if len(split) != 4 {
    err := errors.New("Incorrect number of fields in qrel string: " +line)
    return nil, err
  }

  topic, err := strconv.ParseInt(split[0],10,0)
  if err != nil {
    return nil, err
  }
  iteration := split[1]
  docId := split[2]

  score, err := strconv.ParseInt(split[3],10,0)
  if err != nil {
    return nil, err
  }
  return &Qrel{topic,iteration,docId,score}, nil
}

// This function returns a ResultsFile object created from the
// provided reader (eg a file). 
//
// On errors, a ResultFile containing every Result and ResultList read before the error was encountered is
// returned, along with the error.
func QrelsFromReader(file io.Reader) (QrelsFile,error) {
  var qf QrelsFile
  qf.Qrels = make(map[int64]Qrels)

  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    q, err := QrelFromLine(scanner.Text())
    qrels,ok := qf.Qrels[q.Topic]
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

