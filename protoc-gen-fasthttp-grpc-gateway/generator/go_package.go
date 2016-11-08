package generator

import (
	"fmt"
	"path"
)

type GoPackage struct {
	Path  string
	Name  string
	Alias string
}

func (p GoPackage) String() string {
	if p.Alias == "" || p.Alias == p.Name {
		return fmt.Sprintf("%q", p.Path)
	}

	return fmt.Sprintf("%s %q", p.Alias, p.Path)
}

func NewGoPackage(pkg string) *GoPackage {
	return &GoPackage{
		Path: pkg,
		Name: path.Base(pkg),
	}
}

func NewGoPackageList(list []string) []*GoPackage {
	pkgs := []*GoPackage{}
	pkgMap := map[string]bool{}
	aliasMap := map[string]int{}

	// Remove duplicated packages
	for _, name := range list {
		if pkgMap[name] {
			continue
		}

		pkgMap[name] = true
		pkgs = append(pkgs, NewGoPackage(name))
	}

	// Set alias for packages
	for _, pkg := range pkgs {
		i, ok := aliasMap[pkg.Name]

		if !ok {
			aliasMap[pkg.Name] = 0
			pkg.Alias = pkg.Name
			continue
		}

		aliasMap[pkg.Name] = i + 1
		pkg.Alias = fmt.Sprintf("%s_%d", pkg.Name, i)
	}

	return pkgs
}
