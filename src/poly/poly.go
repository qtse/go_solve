package poly

import (
    "fmt"
    "log"
    "math"
    )

const debug bool = false

/*****************************
 * package config
 ****************************/
func GetMaxDegree() int {
  return max_degree
}

func SetMaxDegree(deg int) {
  max_degree = deg
}

func GetMinCoeff() float64 {
  return min_coeff
}

func SetMinCoeff(deg float64 ) {
  min_coeff = deg
}

func Done() {
  close(multChan)
}

/*****************************
 * func acting on Polynomial
 ****************************/
func NewConstant(c float64) Polynomial {
  if r,found := polyCache[c]; found {
    if debug {
      setLogDefault()
      log.Printf("**DEBUG - Retrieve constant (%f) from cache: %s", c, r.String())
    }
    return r
  }

  if debug {
    setLogDefault()
    log.Printf("**DEBUG - Create new Constant: %f",c)
  }

  r := make(Polynomial)
  if c != 0 {
    r[""] = Term{make(map[string]int), c}
  }
  polyCache[c]=r

  return r
}

func NewPronumeral(p string) Polynomial {
  if r,found := polyCache[p]; found {
    if debug {
      setLogDefault()
      log.Printf("**DEBUG - Retrieve pronumeral (%s) from cache: %s", p, r.String())
    }
    return r
  }

  if debug {
    setLogDefault()
    log.Printf("**DEBUG - Create new pronumeral: %s",p)
  }

  r := make(Polynomial)
  t := Term{make(map[string]int),1}
  t.Pron[p] = 1
  r[t.indexString()] = t
  polyCache[p] = r

  return r
}

func AddPoly(p, q Polynomial) Polynomial {
  if debug {
    setLogDefault()
    log.Printf("**DEBUG - Entering AddPoly: (%s) + (%s)",p.String(),q.String())
  }
  r := p.Copy()
  r.AddPoly(q)

  if debug {
    log.Printf("**DEBUG - (%s) + (%s) = %s",p.String(),q.String(), r.String())
  }
  return r
}

func MultPoly(p,q Polynomial) Polynomial {
  if debug {
    setLogDefault()
    log.Printf("**DEBUG - Entering MultPoly: (%s)*(%s)",p.String(),q.String())
  }

  chs := make([]chan Polynomial, 0, len(p))

  for _,tm := range p {
    rch := make(chan Polynomial)
    chs = append(chs,rch)
    multChan <- &multArg{tm, q, rch}
  }

  r := make(Polynomial)
  for _,rch := range chs {
    r.AddPoly(<-rch)
  }

  if debug {
    log.Printf("**DEBUG - (%s) * (%s) = %s",p.String(), q.String(), r.String())
  }
  return r
}

/*****************************
 * Polynomial class 
 ****************************/
type Polynomial map[string]Term

func (p Polynomial) Copy() Polynomial {
  if debug {
    setLogDefault()
    log.Printf("**DEBUG - Entering Poly.Copy")
  }
  r := make(Polynomial)
  for tstr,term := range p {
    r[tstr] = term.Copy()
  }

  return r
}

func (p Polynomial) AddPoly(q Polynomial) Polynomial {
  if debug {
    setLogDefault()
    log.Printf("**DEBUG - Entering Poly.AddPoly: (%s)+(%s)",p.String(),q.String())
  }
  for tstr, tm := range q {
    _,found := p[tstr]
    if found {
      p[tstr] = Term{tm.Pron, tm.Coeff + p[tstr].Coeff}
    } else {
      p[tstr] = tm.Copy()
    }

    if math.Fabs(p[tstr].Coeff) < min_coeff {
      p[tstr] = tm, false
    }
  }

  return p
}

func (p Polynomial) String() string {
  first := true
  res := ""

  if len(p) == 0 {
    return "0"
  }

  for tstr,term := range p {
    if term.Coeff == 0 {
      p[tstr] = term,false
      continue
    }

    ts,pos := term.absString()
    if !pos {
      res += " - "
    } else if !first {
      res += " + "
    }

    res += ts
    first = false
  }

  return res
}

/*****************************
 * Multiplication code 
 ****************************/
type multArg struct {
  t Term
  q Polynomial
  ch chan Polynomial
}

func multRoutine(ch chan *multArg) {
  for arg := range ch {
    multPoly (arg.t, arg.q, arg.ch)
  }
}

func multPoly (t Term, q Polynomial, ch chan Polynomial) {
  if debug {
    setLogDefault()
    log.Printf("**DEBUG - Entering multPoly: (%s) * (%s)",t.String(),q.String())
  }

  r := make(Polynomial)
  for _,tm := range q {
    nt := Term{make(map[string]int),0}

    for pron,exp := range t.Pron {
      if _,found := tm.Pron[pron]; found {
        nt.Pron[pron] = tm.Pron[pron] + exp
      } else {
        nt.Pron[pron] = exp
      }
    }
    for pron,exp := range tm.Pron {
      if _,found := nt.Pron[pron]; !found {
        nt.Pron[pron] = exp
      }
    }
    if nt.IsDegreeLower(max_degree) {
      nt.Coeff = tm.Coeff * t.Coeff
      r[nt.indexString()] = nt
    }
  }

  if debug {
    log.Printf("**DEBUG - (%s) * (%s) = %s",t.String(), q.String(),r.String())
  }
  ch <- r
  close(ch)
}

/*****************************
 * Term class 
 ****************************/
type Term struct {
  Pron map[string]int
  Coeff float64
}

func (t *Term) indexString() string {
  res := ""
  for p,exp := range (*t).Pron {
    switch exp {
    case 0:
      (*t).Pron[p] = 0,false
      continue
    case 1:
      res += p
    default:
      res += p
      res += "^"
      res += fmt.Sprint(exp)
    }
  }

  return res
}

func (t *Term) absString() (res string,positive bool) {
  if (*t).Coeff == 0 {
    return "",true
  } else if (*t).Coeff > 0 {
    positive = true
  } else {
    positive = false
  }

  if len((*t).Pron) == 0 || math.Fabs((*t).Coeff) != 1 {
    res = fmt.Sprint(math.Fabs((*t).Coeff))
  }

  res += (*t).indexString()

  return
}

func (t *Term) String() string {
  s,p := (*t).absString()
  if !p {
    s = "-" + s
  }
  return s
}

func (t *Term) Copy() Term {
  r := Term{make(map[string]int),(*t).Coeff}
  for p,exp := range (*t).Pron {
    r.Pron[p] = exp
  }
  return r
}

func (t *Term) IsDegreeLower(deg int) bool {
  if deg < 0 {
    return true
  }

  for _,exp := range (*t).Pron {
    deg -= exp
    if deg < 0 {
      return false
    }
  }
  return true
}

func (t *Term) Degree() (deg int) {
  for _,exp := range (*t).Pron {
    deg += exp
  }
  return
}

/*****************************
 * Misc internal functions
 ****************************/
func setLogDefault() {
  log.SetFlags(log.Flags() & ^(log.Ldate|log.Ltime))
}

func init() {
  var i uint
  for i = 0; i < ncpu; i++ {
    go multRoutine(multChan)
  }
}

var multChan chan *multArg = make(chan *multArg, 2*ncpu)
const ncpu uint = 3
var max_degree int = -1
var min_coeff float64 = -1
var polyCache map[interface{}]Polynomial = make(map[interface{}]Polynomial)
