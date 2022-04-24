Alternative to reddits hotness score.
- instead of computing scores in floats, do it in int
- just get the unix now - time of the entity
- older items will have -tive score, current today will have 0
- alternative, just use unix ts

if we want to constrain the date to say, last 30 days
- max(30 - (unix now - date), 0)
- max score is 30 today, 0 after 30 days
- avoid taking all time score, becomes hotness/trending indicator
