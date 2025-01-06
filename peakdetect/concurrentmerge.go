// Copyright (c) 2024 Andrei Gill. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package peakdetect

/*
 Sketch of a concurrent merge operation (very much a WIP; still working on the concept...)

 Basic idea is that we split the stream into 'N' workers, each worker taking two buckets of 3 samples
 each, so, 2 * 3 = 6 samples per worker, so, we divide the whole stream, well, an arbitrary slice/array,
 or as I call it here, a cluster of samples, of the stream, by 6, to get the 'N' workers.

 num_clusters = 2
 samples_per_cluster = 3

 N = len(stream_slice) / (num_clusters * samples_per_cluster)
 N = len(stream_slice) / 6

 Each worker then, computes the 3 samples and merges them, producing a 6-sample cluster.
 When two nearby 6-sample clusters are available, they can be merged into one 12-sample cluster.
 When two nearby 12-sample clusters are available, they can be merged into one 24-sample cluster,
 and so on, and so forth, until only two clusters remain, in which case they're merged, or only
 a single cluster remains, in which case it is the result.

 At the lowest level, we're always going to have 'N' workers. At the next level up, where the
 6-sample clusters are merged... how many workers are we going to have? Half of 'N', i.e., 'N/2'?
 I need to think about this..

 UPDATE:

 At the lowest level of 'N' we run the 3-sample clusters. Then we merge all the workers.
 At the next level, when we merged all the 3-sample clusters, we have 6-sample clusters.

 Now, our:

 samples_per_cluster = 6
 num_clusters = 2

 N_1 = len(stream_slice) / (2 * 6)
 N_1 = len(stream_slice) / 12

 N_2 = len(stream_slice) / (2 * 12)
 N_2 = len(stream_slice) / 24

 N_m = len(stream_slice) / (2 * ..)
 N_m = len(stream_slice) / ..


*/
