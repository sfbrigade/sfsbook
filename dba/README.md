This directory contains the code needed to support the
database that backs the sfsbook application. 

# Details
sfsbook is primarily a document-based storage scheme. So we
use [bleve](http://www.blevesearch.com/) to provide a document
storage scheme above a KV store.

## Refguide
The reference guide itself is a bundle of documents. We pour this into
bleve using an appropriate document mapping that handles the fields
found in the JSON. 

*NB: there are additional details about how to handle schema changes that need
discussion with others*

I am presuming that I can add un-tokenized fields. 

## Comments
An important part of the sfsbook application is the ability of signed-in users to
annotate entries in the book with personal and shared commentary. So, I need a
way to index comments but limit the search results to only those that would
be seen if the user was signed in. 

Idea:

*  treat each comment (per user, per book-entry) as a separate document
*  each comment has a viewability field indicating (specific user, signed-in user, visible to all)
*  then can do search where I require a match
*  comment also needs an owner field. 
*  I'm modeling this loosely after the UNIX permission scheme: me, signed-in, world.
*  Conceivably, more complicated rights scheme could be added.
*  One I do the search over comments, I will get a list of book-entry.
*  I will manually "join" the book-entries up.

## Identities
We have a KV store backing the document store. A separate KV store can manage
the password data. Key is the username. Value is the hashed password, user id combo.
The user-id is needed elsewhere for comments.

## Pending Issues
There remain some issues that I should be worrying about.

*  how to handle updates. Are they efficient or do changes force a large amount of work?


# Deployment
At first while our data volume is small, the reference guide will be compiled
into the go program. This way, deployment is extremely simple. This is not
sustainable in the long term. Later, the documents need to be added via
API.

