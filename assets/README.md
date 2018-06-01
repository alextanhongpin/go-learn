# Slice optimization

Store only what you need. Limiting the slice length can save a lot of memory. The image `slice1.png` and `slice2.png` shows the difference.

The code is basically computing the similarity distance between two users. Since we have roughly 4k users, we need to do a double for loop and compute the score, then save it into the slice. But in the end, we only take the top 20 results. 

We can optimize it by sorting the first 40 results, and then slice the first 20 results everytime we hit `2 x the number of matches we want`. This saves a log of storage.
