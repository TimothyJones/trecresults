Package `trecresults`
--------------------

This is a helper package for reading treceval style `.res` files. 

Example usage:


    import (
      "fmt"
      "github.com/TimothyJones/trecresults"
      "log"
      "os"
      "sort"
    )

    func main() {
      file, err := os.Open("example.res")
      if err != nil {
        log.Fatal(err)
      }
      defer file.Close()

      // Read results from a file
      res,err := trecresults.ResultsFromReader(file)
      if err != nil {
        log.Fatal(err)
      }

      for _,r := range res {
        // print out the result in treceval format
        fmt.Println(r)
      }

      // Normalise the scores in the range 0-1
      r.NormaliseLinear()
    }
