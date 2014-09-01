package trecresults

import (
  "testing"
  "strconv"
)

func CheckResult(r *Result,topic int64,iteration string,docid string,rank int64,score float64,runname string, line string,t *testing.T) {
  if r.Topic != topic {
    t.Error("Expected topic",topic,"got",r.Topic)
  }
  if r.Iteration != iteration {
    t.Error("Expected iteration",iteration,"got",r.Iteration)
  }
  if r.DocId != docid {
    t.Error("Expected DocId",docid,"got",r.DocId)
  }
  if r.Rank != rank {
    t.Error("Expected rank",rank,"got",r.Rank)
  }
  if r.Score != score {
    t.Error("Expected score",score,"got", r.Score)
  }
  if r.RunName != runname {
    t.Error("Expected runname",runname,"got",r.RunName)
  }
  if r.String() != line {
    t.Error("Expected string representation",line,"but got",r)
  }
}

func TestReadLineIncorrectSize(t *testing.T) {
  r,err := ResultFromLine("401 Q0 LA110990-0013 0 13.7471758025085")
  if err == nil {
    t.Error("Expected error, but got nothing")
  }
  if r != nil {
    t.Error("Expected nil result but got",r)
  }

  r,err = ResultFromLine("402 Q1 document 2 12.028 greatrun anotherfield")
  if err == nil {
    t.Error("Expected error, but got nothing")
  }
  if r != nil {
    t.Error("Expected nil result but got",r)
  }
}

func TestReadLineGood(t *testing.T) {
  line1 := "401 Q0 LA110990-0013 0 13.74717580250855 BB2c1.0"
  line2 := "402 Q1 document 2 12.028 greatrun"

  r,err := ResultFromLine(line1)
  if err != nil {
    t.Error("Expected no error, but got",err)
  }
  CheckResult(r,401,"Q0","LA110990-0013",0,13.74717580250855,"BB2c1.0",line1,t)

  r,err = ResultFromLine(line2)
  if err != nil {
    t.Error("Expected no error, but got",err)
  }
  CheckResult(r,402,"Q1","document",2,12.028,"greatrun",line2,t)
}

func TestReadLineBadTopic(t *testing.T) {
  r,err := ResultFromLine("s401 Q0 LA110990-0013 0 13.74717580250855 BB2c1.0")
  if err != nil {
    switch err := err.(type) {
      case *strconv.NumError:
        if err.Func != "ParseInt" {
          t.Error("Error produced not from parse int: got",err.Func)
        }
      default:
        t.Error("Strconv error wasn't produced correctly: got",err)
    }
  } else {
    t.Error("Expected error but got",err)
  }
  if r != nil {
    t.Error("Expected nil response but got",r)
  }
}

func TestReadLineBadRank(t *testing.T) {
  r,err := ResultFromLine("401 Q0 LA110990-0013 rank1 13.74717580250855 BB2c1.0")
  if err != nil {
    switch err := err.(type) {
      case *strconv.NumError:
        if err.Func != "ParseInt" {
          t.Error("Error produced not from parse int: got",err.Func)
        }
      default:
        t.Error("Strconv error wasn't produced correctly: got",err)
    }
  } else {
    t.Error("Expected error but got",err)
  }
  if r != nil {
    t.Error("Expected nil response but got",r)
  }
}

func TestReadLineBadScore(t *testing.T) {
  r,err := ResultFromLine("401 Q0 LA110990-0013 1 1-3.74717580250855 BB2c1.0")
  if err != nil {
    switch err := err.(type) {
      case *strconv.NumError:
        if err.Func != "ParseFloat" {
          t.Error("Error produced not from parsefloat: got",err.Func)
        }
      default:
        t.Error("Strconv error wasn't produced correctly: got",err)
    }
  } else {
    t.Error("Expected error but got",err)
  }
  if r != nil {
    t.Error("Expected nil response but got",r)
  }
}
