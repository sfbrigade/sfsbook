# Implementation Notes
*Wed Aug 31 19:31:02 PDT 2016*

I think that the code as written is both convoluted and needs some additional structure.
I want to figure out how this should work. We have several different flows. I will enumerate
them and the common structure ought to become clearer.

## Static file
Here, the app is running in `debug` mode if there is a file in appropriate
to serve or production mode otherwise.

Flow is like so:

```python
decode cookie
retrieve user record for uuid from cache
if info not  in cache:
	async -> # request to db for user info
if filepath is not cached:
	if desired URL "filepath" is on disk:
		# This means we're running in Debug.
		async -> # request a string by
			open the "filepath"
			load the "filepath" into a string
	else
		# We are in Production mode. We expect to find compressed doc
		async -> # Retrieve desired filepath from static store (go routine adds to cache)
wait <- user_info,  string
if user_info denies right to view:
	cancel string fetch
	redirect # and return an appropriate error

make a  "Reader" from the string
copy from a "Reader" to the "http.ResponseWriter"
```

## Templated
Should work most the same. But with some important differences.

```python
decode cookie
retrieve user record for uuid from cache
# need to handle the difference here between resource and query.
# i.e. URL re-writing.
if info not  in cache:
	async -> # request to db for user info. Which request

# Can insert caching of entire result here.

if "parsed_template" is cached:
	if desired "filepath" is on disk:
		# This means we're running in Debug.
		async -> # request a string by
			open the "filepath"
			load the "filepath" into a string
			parse the string as a template
			return the parsed template
	else
		# We are in Production mode. We expect to find compressed template
		async -> 
			Retrieve desired filepath from static store
			Parse into a template
			<- # return parsed template here
			Inject parsed template into cache
wait <- user_info

augment context with user_info
async -> # request to db for appropriate query. user_info may modify result

wait <- parsed_template, db
if user_info denies right to view:
	cancel parsed_template fetch
	redirect # and return an appropriate error

execute template with db result and user info.
```

## Generated Content
This is content like SASS or Babel transpilation. The details are 
inside of the 

## Summary
The flow has become apparent. I must restructure the code.

## Caching
I want to minimize the amount of work that I need to do for a given page.
What constitutes the final form of a page? Baked output as served to the
requesting UA. This should be cached. What uniquely identifies this?

*  resource path
*  post arguments
*  user uuid

Wait! Shouldn't I benchmark first? Profile. Look at data? Do the simplest
thing? Maybe it doesn't matter. And will never matter.