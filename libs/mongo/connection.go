/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Co., Ltd..
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

/*
 * Revision History:
 *     Initial: 2017/10/28        Feng Yifei
 */

package mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Connection represents a MongoDB collection.
type Connection struct {
	session    *mgo.Session
	collection *mgo.Collection
	Database   string
	Name       string
}

func (conn *Connection) Collection() *mgo.Collection {
	return conn.collection
}

// NewConnection creates a connection to MongoDB.
func NewConnection(s *mgo.Session, db, cname string) *Connection {
	return &Connection{
		session:  s,
		Database: db,
		Name:     cname,
	}
}

// Connect to MongoDB.
func (conn *Connection) Connect() *Connection {
	s := conn.session.Copy()
	c := s.DB(conn.Database).C(conn.Name)

	return &Connection{
		session:    s,
		collection: c,
	}
}

// Disconnect from MongoDB.
func (conn *Connection) Disconnect() {
	conn.session.Close()
}

// GetByID get a single record by ID
func (conn *Connection) GetByID(id interface{}, i interface{}) error {
	return conn.collection.FindId(id).One(i)
}

// GetUniqueOne get a single record by query
func (conn *Connection) GetUniqueOne(q interface{}, doc interface{}) error {
	return conn.collection.Find(q).One(doc)
}

// GetMany get multiple records based on a condition
func (conn *Connection) GetMany(q interface{}, doc interface{}, fields ...string) error {
	if len(fields) == 0 {
		return conn.collection.Find(q).All(doc)
	}

	return conn.collection.Find(q).Sort(fields...).All(doc)
}

// Insert add new documents to a collection.
func (conn *Connection) Insert(doc interface{}) error {
	return conn.collection.Insert(doc)
}

// UpdateByQueryField modify all eligible documents.
func (conn *Connection) UpdateByQueryField(q interface{}, field string, value interface{}) (*mgo.ChangeInfo, error) {
	return conn.collection.UpdateAll(q, bson.M{"$set": bson.M{field: value}})
}

// Update modify existing documents in a collection.
func (conn *Connection) Update(query interface{}, i interface{}) error {
	return conn.collection.Update(query, i)
}

// Upsert creates a new document and inserts it if no documents match the specified filter.
// If there are matching documents, then the operation modifies or replaces the matching document or documents.
func (conn *Connection) Upsert(query interface{}, i interface{}) (*mgo.ChangeInfo, error) {
	return conn.collection.Upsert(query, i)
}

// Delete remove documents from a collection.
func (conn *Connection) Delete(query interface{}) error {
	return conn.collection.Remove(query)
}

// IterAll prepares a pipeline to aggregate and executes the pipeline, works like Iter.All.
func (conn *Connection) IterAll(pipeline interface{}, i interface{}) error {
	return conn.collection.Pipe(pipeline).All(i)
}

// GetLimitedRecords obtain records based on specified conditions.
// The results of the specified number of returns are sorted by the specified fields.
func (conn *Connection) GetLimitedRecords(q interface{}, n int, doc interface{}, fields ...string) error {
	if len(fields) == 0 {
		return conn.collection.Find(q).Limit(n).All(doc)
	}

	return conn.collection.Find(q).Sort(fields...).Limit(n).All(doc)
}
