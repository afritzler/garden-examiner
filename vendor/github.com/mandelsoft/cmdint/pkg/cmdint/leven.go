package cmdint

import "unicode/utf8"

func Levenshtein(a, b string) int {
	f := make([]int, utf8.RuneCountInString(b)+1)

	for j := range f {
		f[j] = j
	}

	for _, ca := range a {
		j := 1
		fj1 := f[0] // fj1 is the value of f[j - 1] in last iteration
		f[0]++
		for _, cb := range b {
			mn := min(f[j]+1, f[j-1]+1) // delete & insert
			if cb != ca {
				mn = min(mn, fj1+1) // change
			} else {
				mn = min(mn, fj1) // matched
			}

			fj1, f[j] = f[j], mn // save f[j] to fj1(j is about to increase), update f[j] to mn
			j++
		}
	}

	return f[len(f)-1]
}

func min(a, b int) int {
	if a <= b {
		return a
	} else {
		return b
	}
}

func SelectBest(name string, candidates ...string) (string, int) {
	c := ""
	min := -1
	for _, n := range candidates {
		d := Levenshtein(n, name)
		if d < len(n)/2 && (min == -1 || min > d) {
			c, min = n, d
		} else {
			if d == min {
				c = ""
			}
		}
	}
	if c == "" {
		min = -1
		for _, n := range candidates {
			if match_omit(n, name) {
				d := Levenshtein(n, name)
				if min == -1 || min > d {
					c, min = n, d
				} else {
					if d == min {
						c = ""
					}
				}
			}
		}
	}
	return c, min
}

func match_omit(cand, pat string) bool {
	if cand == "" {
		return false
	}
	next, l := utf8.DecodeRuneInString(cand)

	for i, r := range pat {
		for r != next {
			if i == 0 {
				return false
			}
			cand = cand[l:]
			if cand == "" {
				return false
			}
			next, l = utf8.DecodeRuneInString(cand)
		}
	}
	return true
}
