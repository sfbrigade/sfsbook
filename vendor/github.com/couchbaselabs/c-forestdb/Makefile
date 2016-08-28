LIB_SOURCES =                             \
          src/api_wrapper.cc        			\
          src/avltree.cc            			\
          src/bgflusher.cc          			\
          src/blockcache.cc         			\
          src/btree.cc              			\
          src/btree_kv.cc           			\
          src/btree_str_kv.cc       			\
          src/btree_fast_str_kv.cc  			\
          src/btreeblock.cc         			\
          src/checksum.cc           			\
          src/compactor.cc          			\
          src/configuration.cc      			\
          src/docio.cc              			\
          src/encryption.cc         			\
          src/encryption_aes.cc     			\
          src/encryption_bogus.cc   			\
          src/fdb_errors.cc         			\
          src/filemgr.cc            			\
          src/filemgr_ops.cc        			\
					src/filemgr_ops_linux.cc        \
					src/filemgr_ops_windows.cc      \
          src/forestdb.cc           			\
          src/hash.cc               			\
          src/hash_functions.cc     			\
          src/hbtrie.cc             			\
          src/iterator.cc           			\
          src/kv_instance.cc        			\
          src/list.cc               			\
          src/snapshot.cc           			\
          src/staleblock.cc         			\
          src/superblock.cc         			\
          src/transaction.cc        			\
          src/version.cc            			\
          src/wal.cc                			\
					utils/debug.cc									\
					utils/system_resource_stats.cc	\
					utils/time_utils.cc							\
					utils/crc32.cc									\
					utils/partiallock.cc						\

sources:
						@for source_file in $(LIB_SOURCES); do \
							echo internal/$$source_file; \
						done
