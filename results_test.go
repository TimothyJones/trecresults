package trecresults

import (
  "testing"
  "strconv"
  "strings"
)


// Helper function to check that a result is as expected, and that it correctly reparses as a string
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


// Checks that errors are thrown when there are too many or too few
// fields in this result
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

// Checks two different well formed result lines
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

// Checks that the correct error is thrown when a non-integer topic ID is presented
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

// Checks that the correct error is thrown when a non-integer rank is provided
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

// Checks that the correct error is thrown when the score cannot be parsed
// as a float
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

// Checks that a sample results file parses correctly
func TestResultsFromFile(t *testing.T) {
  results, err := ResultsFromReader(strings.NewReader(`401 Q0 LA110990-0013 0 13.74717580250855 BB2c1.0
401 Q0 FBIS3-18833 1 13.662447072667604 BB2c1.0
401 Q0 FBIS3-39117 2 13.640016012221363 BB2c1.0
401 Q0 FT941-230 3 13.4799521334611 BB2c1.0
401 Q0 FT924-1346 4 13.418277205894087 BB2c1.0
401 Q0 FT941-4640 5 13.32332784351334 BB2c1.0
401 Q0 LA122190-0057 6 13.278646892401042 BB2c1.0
401 Q0 FBIS3-18916 7 13.00539383125854 BB2c1.0
401 Q0 LA030690-0168 8 12.870710238224662 BB2c1.0
401 Q0 FBIS3-17077 9 12.806848508228754 BB2c1.0
`))
  if err != nil {
    t.Error("Expected no error, but got",err)
  }
  if len(results) != 10 {
    t.Error("Expected 10 results, but got",len(results))
  }

  CheckResult(results[0],401,"Q0","LA110990-0013",0,13.74717580250855,"BB2c1.0","401 Q0 LA110990-0013 0 13.74717580250855 BB2c1.0",t)
  CheckResult(results[1],401,"Q0","FBIS3-18833",1,13.662447072667604,"BB2c1.0","401 Q0 FBIS3-18833 1 13.662447072667604 BB2c1.0",t)
  CheckResult(results[2],401,"Q0","FBIS3-39117",2,13.640016012221363,"BB2c1.0","401 Q0 FBIS3-39117 2 13.640016012221363 BB2c1.0",t)
  CheckResult(results[3],401,"Q0","FT941-230",3,13.4799521334611,"BB2c1.0","401 Q0 FT941-230 3 13.4799521334611 BB2c1.0",t)
  CheckResult(results[4],401,"Q0","FT924-1346",4,13.418277205894087,"BB2c1.0","401 Q0 FT924-1346 4 13.418277205894087 BB2c1.0",t)
  CheckResult(results[5],401,"Q0","FT941-4640",5,13.32332784351334,"BB2c1.0","401 Q0 FT941-4640 5 13.32332784351334 BB2c1.0",t)
  CheckResult(results[6],401,"Q0","LA122190-0057",6,13.278646892401042,"BB2c1.0","401 Q0 LA122190-0057 6 13.278646892401042 BB2c1.0",t)
  CheckResult(results[7],401,"Q0","FBIS3-18916",7,13.00539383125854,"BB2c1.0","401 Q0 FBIS3-18916 7 13.00539383125854 BB2c1.0",t)
  CheckResult(results[8],401,"Q0","LA030690-0168",8,12.870710238224662,"BB2c1.0","401 Q0 LA030690-0168 8 12.870710238224662 BB2c1.0",t)
  CheckResult(results[9],401,"Q0","FBIS3-17077",9,12.806848508228754,"BB2c1.0","401 Q0 FBIS3-17077 9 12.806848508228754 BB2c1.0",t)
}

// Checks that a sample results file correctly normalises
func TestResultsNormaliseLinear(t *testing.T) {
  results, err := ResultsFromReader(strings.NewReader(`401 Q0 LA110990-0013 0 13.74717580250855 BB2c1.0
401 Q0 FBIS3-18833 1 13.662447072667604 BB2c1.0
401 Q0 FBIS3-39117 2 13.640016012221363 BB2c1.0
401 Q0 FT941-230 3 13.4799521334611 BB2c1.0
401 Q0 FT924-1346 4 13.418277205894087 BB2c1.0
401 Q0 FT941-4640 5 13.32332784351334 BB2c1.0
401 Q0 LA122190-0057 6 13.278646892401042 BB2c1.0
401 Q0 FBIS3-18916 7 13.00539383125854 BB2c1.0
401 Q0 LA030690-0168 8 12.870710238224662 BB2c1.0
401 Q0 FBIS3-17077 9 12.806848508228754 BB2c1.0
`))
  if err != nil {
    t.Error("Expected no error, but got",err)
  }
  if len(results) != 10 {
    t.Error("Expected 10 results, but got",len(results))
  }

  results.NormaliseLinear()

  CheckResult(results[0],401,"Q0","LA110990-0013",0,1,"BB2c1.0","401 Q0 LA110990-0013 0 1 BB2c1.0",t)
  CheckResult(results[1],401,"Q0","FBIS3-18833",1,0.9098944268061034,"BB2c1.0","401 Q0 FBIS3-18833 1 0.9098944268061034 BB2c1.0",t)
  CheckResult(results[2],401,"Q0","FBIS3-39117",2,0.8860399023413839,"BB2c1.0","401 Q0 FBIS3-39117 2 0.8860399023413839 BB2c1.0",t)
  CheckResult(results[3],401,"Q0","FT941-230",3,0.7158184488815479,"BB2c1.0","401 Q0 FT941-230 3 0.7158184488815479 BB2c1.0",t)
  CheckResult(results[4],401,"Q0","FT924-1346",4,0.6502296608689117,"BB2c1.0","401 Q0 FT924-1346 4 0.6502296608689117 BB2c1.0",t)
  CheckResult(results[5],401,"Q0","FT941-4640",5,0.5492548588416348,"BB2c1.0","401 Q0 FT941-4640 5 0.5492548588416348 BB2c1.0",t)
  CheckResult(results[6],401,"Q0","LA122190-0057",6,0.5017384766371606,"BB2c1.0","401 Q0 LA122190-0057 6 0.5017384766371606 BB2c1.0",t)
  CheckResult(results[7],401,"Q0","FBIS3-18916",7,0.211144911178883,"BB2c1.0","401 Q0 FBIS3-18916 7 0.211144911178883 BB2c1.0",t)
  CheckResult(results[8],401,"Q0","LA030690-0168",8,0.06791436384373004,"BB2c1.0","401 Q0 LA030690-0168 8 0.06791436384373004 BB2c1.0",t)
  CheckResult(results[9],401,"Q0","FBIS3-17077",9,0,"BB2c1.0","401 Q0 FBIS3-17077 9 0 BB2c1.0",t)
}

