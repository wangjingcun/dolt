// This file was generated by nomdl/codegen.

package gen

import (
	"github.com/attic-labs/noms/ref"
	"github.com/attic-labs/noms/types"
)

var __genPackageInFile_struct_with_imports_CachedRef ref.Ref

// This function builds up a Noms value that describes the type
// package implemented by this file and registers it with the global
// type package definition cache.
func init() {
	p := types.NewPackage([]types.TypeRef{
		types.MakeEnumTypeRef("LocalE", "LocalE1", "Ignored"),
		types.MakeStructTypeRef("ImportUser",
			[]types.Field{
				types.Field{"importedStruct", types.MakeTypeRef(ref.Parse("sha1-09d2fdd9743c4daec6deebbbc1a38f75ad088eca"), 0), false},
				types.Field{"enum", types.MakeTypeRef(ref.Ref{}, 0), false},
			},
			types.Choices{},
		),
	}, []ref.Ref{
		ref.Parse("sha1-09d2fdd9743c4daec6deebbbc1a38f75ad088eca"),
	})
	__genPackageInFile_struct_with_imports_CachedRef = types.RegisterPackage(&p)
}

// LocalE

type LocalE uint32

const (
	LocalE1 LocalE = iota
	Ignored
)

func NewLocalE() LocalE {
	return LocalE(0)
}

var __typeRefForLocalE types.TypeRef

func (e LocalE) TypeRef() types.TypeRef {
	return __typeRefForLocalE
}

func init() {
	__typeRefForLocalE = types.MakeTypeRef(__genPackageInFile_struct_with_imports_CachedRef, 0)
	types.RegisterFromValFunction(__typeRefForLocalE, func(v types.Value) types.Value {
		return LocalE(uint32(v.(types.UInt32)))
	})
}

func (e LocalE) InternalImplementation() uint32 {
	return uint32(e)
}

func (e LocalE) Equals(other types.Value) bool {
	return e == other
}

func (e LocalE) Ref() ref.Ref {
	throwaway := ref.Ref{}
	return types.EnsureRef(&throwaway, e)
}

func (e LocalE) Chunks() []ref.Ref {
	return nil
}

// ImportUser

type ImportUser struct {
	_importedStruct D
	_enum           LocalE

	ref *ref.Ref
}

func NewImportUser() ImportUser {
	return ImportUser{
		_importedStruct: NewD(),
		_enum:           NewLocalE(),

		ref: &ref.Ref{},
	}
}

type ImportUserDef struct {
	ImportedStruct DDef
	Enum           LocalE
}

func (def ImportUserDef) New() ImportUser {
	return ImportUser{
		_importedStruct: def.ImportedStruct.New(),
		_enum:           def.Enum,
		ref:             &ref.Ref{},
	}
}

func (s ImportUser) Def() (d ImportUserDef) {
	d.ImportedStruct = s._importedStruct.Def()
	d.Enum = s._enum
	return
}

var __typeRefForImportUser types.TypeRef
var __typeDefForImportUser types.TypeRef

func (m ImportUser) TypeRef() types.TypeRef {
	return __typeRefForImportUser
}

func init() {
	__typeRefForImportUser = types.MakeTypeRef(__genPackageInFile_struct_with_imports_CachedRef, 1)
	__typeDefForImportUser = types.MakeStructTypeRef("ImportUser",
		[]types.Field{
			types.Field{"importedStruct", types.MakeTypeRef(ref.Parse("sha1-09d2fdd9743c4daec6deebbbc1a38f75ad088eca"), 0), false},
			types.Field{"enum", types.MakeTypeRef(__genPackageInFile_struct_with_imports_CachedRef, 0), false},
		},
		types.Choices{},
	)
	types.RegisterStructBuilder(__typeRefForImportUser, builderForImportUser)
}

func (s ImportUser) InternalImplementation() types.Struct {
	// TODO: Remove this
	m := map[string]types.Value{
		"importedStruct": s._importedStruct,
		"enum":           s._enum,
	}
	return types.NewStruct(__typeRefForImportUser, __typeDefForImportUser, m)
}

func builderForImportUser() chan types.Value {
	c := make(chan types.Value)
	s := ImportUser{ref: &ref.Ref{}}
	go func() {
		s._importedStruct = (<-c).(D)
		s._enum = (<-c).(LocalE)

		c <- s
	}()
	return c
}

func (s ImportUser) Equals(other types.Value) bool {
	return other != nil && __typeRefForImportUser.Equals(other.TypeRef()) && s.Ref() == other.Ref()
}

func (s ImportUser) Ref() ref.Ref {
	return types.EnsureRef(s.ref, s)
}

func (s ImportUser) Chunks() (chunks []ref.Ref) {
	chunks = append(chunks, __typeRefForImportUser.Chunks()...)
	chunks = append(chunks, s._importedStruct.Chunks()...)
	return
}

func (s ImportUser) ImportedStruct() D {
	return s._importedStruct
}

func (s ImportUser) SetImportedStruct(val D) ImportUser {
	s._importedStruct = val
	s.ref = &ref.Ref{}
	return s
}

func (s ImportUser) Enum() LocalE {
	return s._enum
}

func (s ImportUser) SetEnum(val LocalE) ImportUser {
	s._enum = val
	s.ref = &ref.Ref{}
	return s
}

// ListOfD

type ListOfD struct {
	l   types.List
	ref *ref.Ref
}

func NewListOfD() ListOfD {
	return ListOfD{types.NewList(), &ref.Ref{}}
}

type ListOfDDef []DDef

func (def ListOfDDef) New() ListOfD {
	l := make([]types.Value, len(def))
	for i, d := range def {
		l[i] = d.New()
	}
	return ListOfD{types.NewList(l...), &ref.Ref{}}
}

func (l ListOfD) Def() ListOfDDef {
	d := make([]DDef, l.Len())
	for i := uint64(0); i < l.Len(); i++ {
		d[i] = l.l.Get(i).(D).Def()
	}
	return d
}

func (l ListOfD) InternalImplementation() types.List {
	return l.l
}

func (l ListOfD) Equals(other types.Value) bool {
	return other != nil && __typeRefForListOfD.Equals(other.TypeRef()) && l.Ref() == other.Ref()
}

func (l ListOfD) Ref() ref.Ref {
	return types.EnsureRef(l.ref, l)
}

func (l ListOfD) Chunks() (chunks []ref.Ref) {
	chunks = append(chunks, l.TypeRef().Chunks()...)
	chunks = append(chunks, l.l.Chunks()...)
	return
}

// A Noms Value that describes ListOfD.
var __typeRefForListOfD types.TypeRef

func (m ListOfD) TypeRef() types.TypeRef {
	return __typeRefForListOfD
}

func init() {
	__typeRefForListOfD = types.MakeCompoundTypeRef(types.ListKind, types.MakeTypeRef(ref.Parse("sha1-09d2fdd9743c4daec6deebbbc1a38f75ad088eca"), 0))
	types.RegisterFromValFunction(__typeRefForListOfD, func(v types.Value) types.Value {
		return ListOfD{v.(types.List), &ref.Ref{}}
	})
}

func (l ListOfD) Len() uint64 {
	return l.l.Len()
}

func (l ListOfD) Empty() bool {
	return l.Len() == uint64(0)
}

func (l ListOfD) Get(i uint64) D {
	return l.l.Get(i).(D)
}

func (l ListOfD) Slice(idx uint64, end uint64) ListOfD {
	return ListOfD{l.l.Slice(idx, end), &ref.Ref{}}
}

func (l ListOfD) Set(i uint64, val D) ListOfD {
	return ListOfD{l.l.Set(i, val), &ref.Ref{}}
}

func (l ListOfD) Append(v ...D) ListOfD {
	return ListOfD{l.l.Append(l.fromElemSlice(v)...), &ref.Ref{}}
}

func (l ListOfD) Insert(idx uint64, v ...D) ListOfD {
	return ListOfD{l.l.Insert(idx, l.fromElemSlice(v)...), &ref.Ref{}}
}

func (l ListOfD) Remove(idx uint64, end uint64) ListOfD {
	return ListOfD{l.l.Remove(idx, end), &ref.Ref{}}
}

func (l ListOfD) RemoveAt(idx uint64) ListOfD {
	return ListOfD{(l.l.RemoveAt(idx)), &ref.Ref{}}
}

func (l ListOfD) fromElemSlice(p []D) []types.Value {
	r := make([]types.Value, len(p))
	for i, v := range p {
		r[i] = v
	}
	return r
}

type ListOfDIterCallback func(v D, i uint64) (stop bool)

func (l ListOfD) Iter(cb ListOfDIterCallback) {
	l.l.Iter(func(v types.Value, i uint64) bool {
		return cb(v.(D), i)
	})
}

type ListOfDIterAllCallback func(v D, i uint64)

func (l ListOfD) IterAll(cb ListOfDIterAllCallback) {
	l.l.IterAll(func(v types.Value, i uint64) {
		cb(v.(D), i)
	})
}

type ListOfDFilterCallback func(v D, i uint64) (keep bool)

func (l ListOfD) Filter(cb ListOfDFilterCallback) ListOfD {
	nl := NewListOfD()
	l.IterAll(func(v D, i uint64) {
		if cb(v, i) {
			nl = nl.Append(v)
		}
	})
	return nl
}
