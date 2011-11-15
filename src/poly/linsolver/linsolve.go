package linsolver

import (
    "fmt"
    "os"
    "poly"
    )

const ncpu int = 3

type Solution map[string]float64

func Done() {
  close(subChan)
}

func (s Solution) String() string {
  res := ""
  first := true
  for p,v := range s {
    if !first {
      res += ", "
    } else {
      first = false
    }

    res += fmt.Sprintf("%s:%f",p,v)
  }
  return res
}

func Solve(system []poly.Polynomial) (Solution, os.Error) {
  m,symTbl,err := formMatrix(system)
  if err != nil {
    return nil,err
  }
  todo := make(map[string]int)
  for p,i := range symTbl {
    todo[p] = i
  }

  for ri,r := range m {
    vIdx := -1
    if v,_,ok,z := r.solved(todo);z {
      continue
    } else if ok {
      vIdx = todo[v]
    } else {
      for i,c  := range r[:len(r)-1] {
        if c != 0 {
          vIdx = i
          break
        }
      }
    }

    r.normalise(1/r[vIdx])
    req := 0
    retChan := make(chan *subRes,len(m))
    for ri2,r2 := range m {
      if ri2 == ri {
        continue
      }
      if _,_,ok,_ := r2.solved(todo);ok {
        continue
      }
      subChan <- &subArg{r2,r,r2[vIdx],retChan}
      req++
    }
    for i := 0; i < req; i++ {
      <-retChan
    }
  }

  return m.getSolution(symTbl)
}

func formMatrix(system []poly.Polynomial) (m matrix, symTbl map[string]int, err os.Error) {
  symTbl = make(map[string]int)
  idx := 0
  for _,p := range system {
    for _,t := range p {
      for pron,_ := range t.Pron {
        if _,found := symTbl[pron]; !found {
          symTbl[pron]=idx
          idx++
        }
      }
    }
  }

  m = make(matrix,len(system),len(system))
  for i,p := range system {
    m[i] = make(row,len(symTbl)+1,len(symTbl)+1)

    for _,t := range p {
      switch t.Degree() {
      case 0:
        idx = len(symTbl)
      case 1:
        for pron,_ := range t.Pron {
          // This should only be executed once
          var ok bool
          if idx,ok = symTbl[pron]; !ok {
            err = os.NewError("internal error - pronumeral not found")
            return
          }
        }
      default:
        err = os.NewError("Non linear equation")
        return
      }
      m[i][idx] = t.Coeff
    }
  }
  return
}

func (r row) subtract(ref row, coeff float64) (res row, err os.Error) {
  if len(r) != len(ref) {
    return nil,os.NewError("row length mismatch")
  }
  for i := range r {
    r[i] -= coeff * ref[i]
  }
  return r,nil
}

func (r row) normalise(scale float64) row {
  for i := range r {
    r[i] *= scale
  }
  return r
}

func (r row) solved(todoTbl map[string]int) (string,float64,bool,bool) {
  pStr := ""
  for v,i := range todoTbl {
    if pStr != "" && r[i] != 0 {
      return "",0,false,false
    } else if r[i] != 0 {
      pStr = v
    }
  }
  if pStr == "" {
    return pStr,r[len(r)-1],true,true
  }
  return pStr,r[len(r)-1],true,false
}

func (r row) String() (res string) {
  for _,c := range r {
    res += fmt.Sprintf("%f ",c)
  }
  res += fmt.Sprintln()
  return
}

func (m matrix) String() (res string) {
  for _,r := range m {
    res += r.String()
  }
  return
}

func (m matrix) getSolution(symTbl map[string]int) (res Solution, err os.Error) {
  res = make(Solution)
  for _,r := range m {
    if p,v,ok,z := r.solved(symTbl); !z {
      if ok {
        if vv,found := res[p]; found && v!=vv {
          err = os.NewError(fmt.Sprintf("System inconsistent: %f != %f",v,vv))
          res = nil
          return
        }
        res[p] = -v
      } else {
        err = os.NewError("Not completely solved")
      }
    }
  }

  return
}

var subChan chan *subArg = make(chan *subArg,3)
func goSub(arg chan *subArg) {
  for a := range subChan {
    arg := *a
    r,e := arg.From.subtract(arg.Ref, arg.Coeff)
    arg.Return <- &subRes{r,e}
  }
}

func init() {
  for i := 0; i < ncpu; i++ {
    go goSub(subChan)
  }
}

type subArg struct {
  From row
  Ref row
  Coeff float64
  Return chan *subRes
}
type subRes struct {
  res row
  err os.Error
}

type row []float64
type matrix []row
