// Package trecresults provides helper functions for reading and writing trec results files
// suitable for using with treceval
package trecresults

import (
  "errors"
  "fmt"
  "strconv"
  "strings"
  "io"
  "bufio"
)

// Describes a single entry in a trec result list
type Result struct {
  Topic int64      // the integer topic ID
  Iteration string // the iteration this run is associated with (ignored by treceval)
  DocId string     // the document ID for this result
  Rank int64       // the rank in the result list
  Score float64    // the score the document received for this topic
  RunName string   // the name of the run this result is from
}

// Type definition for a result list
type ResultList []*Result

// Formats a result structure into the original string representation that can be used with treceval
func (r *Result) String() string {
  return fmt.Sprintf("%d %s %s %d %g %s",r.Topic,r.Iteration,r.DocId,r.Rank,r.Score,r.RunName)
}


// Creates a result structure from a single line from a results file
// Returns parsing errors if any of the integer or float fields do not parse
// Returns an error if there are not 6 fields in the result line
// On error, a nil result is returned
func ResultFromLine(line string) (*Result, error) {
  split := strings.Fields(line)

  if len(split) != 6 {
    err := errors.New("Incorrect number of fields in result string: " +line)
    return nil, err
  }

  topic, err := strconv.ParseInt(split[0],10,0)
  if err != nil {
    return nil, err
  }
  iteration := split[1]
  docId := split[2]

  rank, err := strconv.ParseInt(split[3],10,0)
  if err != nil {
    return nil, err
  }

  score, err := strconv.ParseFloat(split[4],64)
  if err != nil {
    return nil, err
  }
  runname := split[5]

  return &Result{topic,iteration,docId,rank,score,runname}, nil
}

// This function returns a silce of results created from the
// provided reader (eg a file). 
// On errors,a slice of every result read up until the error was encountered is
// returned, along with the error
func ResultsFromReader(file io.Reader) (ResultList,error) {
  results := make([]*Result,0,0)

  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    r, err := ResultFromLine(scanner.Text())
    if err != nil {
      return results, err
    }
    results = append(results,r)
  }

  if err := scanner.Err(); err != nil {
    return results, err
  }
  return results, nil
}

// This function operates on a slice of results, and normalises the score
// of each result by score (score - min)/(max - min). This puts scores
// in to the range 0-1, where 1 is the highest score, and 0 is the lowest.
//
// No assumptions are made as to whether the slice is pre sorted
func (r ResultList) NormaliseLinear() {
  if len(r) == 0 {
    return
  }
  max := r[0].Score
  min := r[0].Score
  for _,res := range r {
    if res.Score > max {
      max = res.Score
    }
    if res.Score < min {
      min = res.Score
    }
  }

  for _,res := range r {
    res.Score = (res.Score - min)/(max - min)
  }
}
