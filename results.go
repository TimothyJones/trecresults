package trecresults

import (
  "errors"
  "fmt"
  "strconv"
  "strings"
  "io"
  "bufio"
)

// 401 Q0 LA110990-0013 0 13.74717580250855 BB2c1.0
type Result struct {
  Topic int64
  Iteration string
  DocId string
  Rank int64
  Score float64
  RunName string
}

func (r *Result) String() string {
  return fmt.Sprintf("%d %s %s %d %g %s",r.Topic,r.Iteration,r.DocId,r.Rank,r.Score,r.RunName)
}

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

func ResultsFromReader(file io.Reader) ([]*Result,error) {
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
