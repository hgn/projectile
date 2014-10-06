
# Installation #

- go get github.com/gorilla/sessions
- go get github.com/gorilla/mux

# Concepts #

## Items ##

- Can have a deadline
- Can be assigned to a Project/Task/Deadline
- Can have tags (item tag namespace)
- Can have one or more responsibles
- Can have a priority
- Cannot have a duration

In shprt, items are notetracker and are often the first step in forming a task
or deadline. But this is optional. Normally items are just a puzzle piece in on
larger task. Items provide place to assign important inforatmtion to tasks or
deadlines. Often items are written directly in a meeting. After creating a item
it is active. THe only other state is archieved.

Items can later be converted to Tasks or Deadlines - but not to Projects.

## Project ##

## Tasks ##

## Deadlines ##


# REST API #

/api/users
{
	"users": [ "user1", "user2" ]
}

/api/user/<user>

{
	"nick" : "<nick>"
}
