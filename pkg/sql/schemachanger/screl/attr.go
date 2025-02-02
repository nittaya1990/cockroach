// Copyright 2021 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package screl

import (
	"reflect"

	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/rel"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scpb"
	"github.com/cockroachdb/cockroach/pkg/util/protoutil"
)

// Attr are attributes used to identify, order, and relate nodes,
// targets, and elements to each other using the rel library.
type Attr int

// MustQuery constructs a query using this package's schema. Intending to be
// called during init, this function panics if query construction fails.
func MustQuery(clauses ...rel.Clause) *rel.Query {
	q, err := rel.NewQuery(Schema, clauses...)
	if err != nil {
		panic(err)
	}
	return q
}

var _ rel.Attr = Attr(0)

//go:generate stringer -type=Attr -trimprefix=Attr
const (
	_ Attr = iota // reserve 0 for rel.Type
	// DescID is the descriptor ID to which this element belongs.
	DescID
	// ReferencedDescID is the descriptor ID to which this element refers.
	ReferencedDescID
	//ColumnID is the column ID to which this element corresponds.
	ColumnID
	// Name is the name of the element.
	Name
	// IndexID is the index ID to which this element corresponds.
	IndexID
	// Direction is the direction of a Target.
	Direction
	// Status is the Status of a Node.
	Status
	// Element references an element.
	Element
	// Target is the reference from a node to a target.
	Target
)

var t = reflect.TypeOf

// Schema is the schema exported by this package covering the elements of scpb.
var Schema = rel.MustSchema("screl",
	rel.AttrType(Element, t((*protoutil.Message)(nil)).Elem()),
	rel.EntityMapping(
		t((*scpb.Node)(nil)),
		rel.EntityAttr(Status, "Status"),
		rel.EntityAttr(Target, "Target"),
	),
	rel.EntityMapping(
		t((*scpb.Target)(nil)),
		rel.EntityAttr(Direction, "Direction"),
		rel.EntityAttr(Element,
			"Column",
			"PrimaryIndex",
			"SecondaryIndex",
			"SequenceDependency",
			"UniqueConstraint",
			"CheckConstraint",
			"Sequence",
			"DefaultExpression",
			"View",
			"TypeRef",
			"Table",
			"OutForeignKey",
			"InForeignKey",
			"RelationDependedOnBy",
			"SequenceOwner",
			"Type",
			"Schema",
			"Database",
		),
	),
	rel.EntityMapping(t((*scpb.Column)(nil)),
		rel.EntityAttr(DescID, "TableID"),
		rel.EntityAttr(ColumnID, "Column.ID"),
		rel.EntityAttr(Name, "Column.Name"),
	),
	rel.EntityMapping(t((*scpb.PrimaryIndex)(nil)),
		rel.EntityAttr(DescID, "TableID"),
		rel.EntityAttr(IndexID, "IndexId"),
		rel.EntityAttr(Name, "IndexName"),
	),
	rel.EntityMapping(t((*scpb.SecondaryIndex)(nil)),
		rel.EntityAttr(DescID, "TableID"),
		rel.EntityAttr(IndexID, "IndexId"),
		rel.EntityAttr(Name, "IndexName"),
	),
	rel.EntityMapping(t((*scpb.SequenceDependency)(nil)),
		rel.EntityAttr(DescID, "SequenceID"),
		rel.EntityAttr(ReferencedDescID, "TableID"),
		rel.EntityAttr(ColumnID, "ColumnID"),
	),
	rel.EntityMapping(t((*scpb.UniqueConstraint)(nil)),
		rel.EntityAttr(DescID, "TableID"),
		rel.EntityAttr(IndexID, "IndexID"),
	),
	rel.EntityMapping(t((*scpb.CheckConstraint)(nil)),
		rel.EntityAttr(DescID, "TableID"),
		rel.EntityAttr(Name, "Name"),
	),
	rel.EntityMapping(t((*scpb.Sequence)(nil)),
		rel.EntityAttr(DescID, "SequenceID"),
	),
	rel.EntityMapping(t((*scpb.DefaultExpression)(nil)),
		rel.EntityAttr(DescID, "TableID"),
		rel.EntityAttr(ColumnID, "ColumnID"),
	),
	rel.EntityMapping(t((*scpb.View)(nil)),
		rel.EntityAttr(DescID, "TableID"),
	),
	rel.EntityMapping(t((*scpb.TypeReference)(nil)),
		rel.EntityAttr(DescID, "DescID"),
		rel.EntityAttr(ReferencedDescID, "TypeID"),
	),
	rel.EntityMapping(t((*scpb.Table)(nil)),
		rel.EntityAttr(DescID, "TableID"),
	),
	rel.EntityMapping(t((*scpb.InboundForeignKey)(nil)),
		rel.EntityAttr(DescID, "OriginID"),
		rel.EntityAttr(ReferencedDescID, "ReferenceID"),
		rel.EntityAttr(Name, "Name"),
	),
	rel.EntityMapping(t((*scpb.OutboundForeignKey)(nil)),
		rel.EntityAttr(DescID, "OriginID"),
		rel.EntityAttr(ReferencedDescID, "ReferenceID"),
		rel.EntityAttr(Name, "Name"),
	),
	rel.EntityMapping(t((*scpb.RelationDependedOnBy)(nil)),
		rel.EntityAttr(DescID, "TableID"),
		rel.EntityAttr(ReferencedDescID, "DependedOnBy"),
	),
	rel.EntityMapping(t((*scpb.SequenceOwnedBy)(nil)),
		rel.EntityAttr(DescID, "SequenceID"),
		rel.EntityAttr(ReferencedDescID, "OwnerTableID"),
	),
	rel.EntityMapping(t((*scpb.Type)(nil)),
		rel.EntityAttr(DescID, "TypeID"),
	),
	rel.EntityMapping(t((*scpb.Schema)(nil)),
		rel.EntityAttr(DescID, "SchemaID"),
	),
	rel.EntityMapping(t((*scpb.Database)(nil)),
		rel.EntityAttr(DescID, "DatabaseID"),
	),
	rel.EntityMapping(t((*scpb.Partitioning)(nil)),
		rel.EntityAttr(DescID, "TableID"),
		rel.EntityAttr(IndexID, "IndexId"),
	),
)

// JoinTargetNode generates a clause that joins the target and node vars
// to the corresponding element.
func JoinTargetNode(element, target, node rel.Var) rel.Clause {
	return rel.And(
		target.Type((*scpb.Target)(nil)),
		target.AttrEqVar(Element, element),
		node.Type((*scpb.Node)(nil)),
		node.AttrEqVar(Target, target),
	)
}
