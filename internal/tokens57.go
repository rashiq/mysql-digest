package internal

// This file maps MySQL 8.0 token IDs to MySQL 5.7 token IDs.
//
// MySQL 5.7 token ranges:
//   Main tokens: 258-909 (652 tokens)
//   Hint tokens: 910-932 (23 tokens)
//   Digest tokens: 928-938

const (
	m57_ABORT_SYM                         = 258
	m57_ACCESSIBLE_SYM                    = 259
	m57_ACCOUNT_SYM                       = 260
	m57_ACTION                            = 261
	m57_ADD                               = 262
	m57_ADDDATE_SYM                       = 263
	m57_AFTER_SYM                         = 264
	m57_AGAINST                           = 265
	m57_AGGREGATE_SYM                     = 266
	m57_ALGORITHM_SYM                     = 267
	m57_ALL                               = 268
	m57_ALTER                             = 269
	m57_ALWAYS_SYM                        = 270
	m57_ANALYSE_SYM                       = 271
	m57_ANALYZE_SYM                       = 272
	m57_AND_AND_SYM                       = 273
	m57_AND_SYM                           = 274
	m57_ANY_SYM                           = 275
	m57_AS                                = 276
	m57_ASC                               = 277
	m57_ASCII_SYM                         = 278
	m57_ASENSITIVE_SYM                    = 279
	m57_AT_SYM                            = 280
	m57_AUTOEXTEND_SIZE_SYM               = 281
	m57_AUTO_INC                          = 282
	m57_AVG_ROW_LENGTH                    = 283
	m57_AVG_SYM                           = 284
	m57_BACKUP_SYM                        = 285
	m57_BEFORE_SYM                        = 286
	m57_BEGIN_SYM                         = 287
	m57_BETWEEN_SYM                       = 288
	m57_BIGINT                            = 289
	m57_BINARY                            = 290
	m57_BINLOG_SYM                        = 291
	m57_BIN_NUM                           = 292
	m57_BIT_AND                           = 293
	m57_BIT_OR                            = 294
	m57_BIT_SYM                           = 295
	m57_BIT_XOR                           = 296
	m57_BLOB_SYM                          = 297
	m57_BLOCK_SYM                         = 298
	m57_BOOLEAN_SYM                       = 299
	m57_BOOL_SYM                          = 300
	m57_BOTH                              = 301
	m57_BTREE_SYM                         = 302
	m57_BY                                = 303
	m57_BYTE_SYM                          = 304
	m57_CACHE_SYM                         = 305
	m57_CALL_SYM                          = 306
	m57_CASCADE                           = 307
	m57_CASCADED                          = 308
	m57_CASE_SYM                          = 309
	m57_CAST_SYM                          = 310
	m57_CATALOG_NAME_SYM                  = 311
	m57_CHAIN_SYM                         = 312
	m57_CHANGE                            = 313
	m57_CHANGED                           = 314
	m57_CHANNEL_SYM                       = 315
	m57_CHARSET                           = 316
	m57_CHAR_SYM                          = 317
	m57_CHECKSUM_SYM                      = 318
	m57_CHECK_SYM                         = 319
	m57_CIPHER_SYM                        = 320
	m57_CLASS_ORIGIN_SYM                  = 321
	m57_CLIENT_SYM                        = 322
	m57_CLOSE_SYM                         = 323
	m57_COALESCE                          = 324
	m57_CODE_SYM                          = 325
	m57_COLLATE_SYM                       = 326
	m57_COLLATION_SYM                     = 327
	m57_COLUMNS                           = 328
	m57_COLUMN_SYM                        = 329
	m57_COLUMN_FORMAT_SYM                 = 330
	m57_COLUMN_NAME_SYM                   = 331
	m57_COMMENT_SYM                       = 332
	m57_COMMITTED_SYM                     = 333
	m57_COMMIT_SYM                        = 334
	m57_COMPACT_SYM                       = 335
	m57_COMPLETION_SYM                    = 336
	m57_COMPRESSED_SYM                    = 337
	m57_COMPRESSION_SYM                   = 338
	m57_ENCRYPTION_SYM                    = 339
	m57_CONCURRENT                        = 340
	m57_CONDITION_SYM                     = 341
	m57_CONNECTION_SYM                    = 342
	m57_CONSISTENT_SYM                    = 343
	m57_CONSTRAINT                        = 344
	m57_CONSTRAINT_CATALOG_SYM            = 345
	m57_CONSTRAINT_NAME_SYM               = 346
	m57_CONSTRAINT_SCHEMA_SYM             = 347
	m57_CONTAINS_SYM                      = 348
	m57_CONTEXT_SYM                       = 349
	m57_CONTINUE_SYM                      = 350
	m57_CONVERT_SYM                       = 351
	m57_COUNT_SYM                         = 352
	m57_CPU_SYM                           = 353
	m57_CREATE                            = 354
	m57_CROSS                             = 355
	m57_CUBE_SYM                          = 356
	m57_CURDATE                           = 357
	m57_CURRENT_SYM                       = 358
	m57_CURRENT_USER                      = 359
	m57_CURSOR_SYM                        = 360
	m57_CURSOR_NAME_SYM                   = 361
	m57_CURTIME                           = 362
	m57_DATABASE                          = 363
	m57_DATABASES                         = 364
	m57_DATAFILE_SYM                      = 365
	m57_DATA_SYM                          = 366
	m57_DATETIME                          = 367
	m57_DATE_ADD_INTERVAL                 = 368
	m57_DATE_SUB_INTERVAL                 = 369
	m57_DATE_SYM                          = 370
	m57_DAY_HOUR_SYM                      = 371
	m57_DAY_MICROSECOND_SYM               = 372
	m57_DAY_MINUTE_SYM                    = 373
	m57_DAY_SECOND_SYM                    = 374
	m57_DAY_SYM                           = 375
	m57_DEALLOCATE_SYM                    = 376
	m57_DECIMAL_NUM                       = 377
	m57_DECIMAL_SYM                       = 378
	m57_DECLARE_SYM                       = 379
	m57_DEFAULT                           = 380
	m57_DEFAULT_AUTH_SYM                  = 381
	m57_DEFINER_SYM                       = 382
	m57_DELAYED_SYM                       = 383
	m57_DELAY_KEY_WRITE_SYM               = 384
	m57_DELETE_SYM                        = 385
	m57_DESC                              = 386
	m57_DESCRIBE                          = 387
	m57_DES_KEY_FILE                      = 388
	m57_DETERMINISTIC_SYM                 = 389
	m57_DIAGNOSTICS_SYM                   = 390
	m57_DIRECTORY_SYM                     = 391
	m57_DISABLE_SYM                       = 392
	m57_DISCARD                           = 393
	m57_DISK_SYM                          = 394
	m57_DISTINCT                          = 395
	m57_DIV_SYM                           = 396
	m57_DOUBLE_SYM                        = 397
	m57_DO_SYM                            = 398
	m57_DROP                              = 399
	m57_DUAL_SYM                          = 400
	m57_DUMPFILE                          = 401
	m57_DUPLICATE_SYM                     = 402
	m57_DYNAMIC_SYM                       = 403
	m57_EACH_SYM                          = 404
	m57_ELSE                              = 405
	m57_ELSEIF_SYM                        = 406
	m57_ENABLE_SYM                        = 407
	m57_ENCLOSED                          = 408
	m57_END                               = 409
	m57_ENDS_SYM                          = 410
	m57_END_OF_INPUT                      = 411
	m57_ENGINES_SYM                       = 412
	m57_ENGINE_SYM                        = 413
	m57_ENUM                              = 414
	m57_EQ                                = 415
	m57_EQUAL_SYM                         = 416
	m57_ERROR_SYM                         = 417
	m57_ERRORS                            = 418
	m57_ESCAPED                           = 419
	m57_ESCAPE_SYM                        = 420
	m57_EVENTS_SYM                        = 421
	m57_EVENT_SYM                         = 422
	m57_EVERY_SYM                         = 423
	m57_EXCHANGE_SYM                      = 424
	m57_EXECUTE_SYM                       = 425
	m57_EXISTS                            = 426
	m57_EXIT_SYM                          = 427
	m57_EXPANSION_SYM                     = 428
	m57_EXPIRE_SYM                        = 429
	m57_EXPORT_SYM                        = 430
	m57_EXTENDED_SYM                      = 431
	m57_EXTENT_SIZE_SYM                   = 432
	m57_EXTRACT_SYM                       = 433
	m57_FALSE_SYM                         = 434
	m57_FAST_SYM                          = 435
	m57_FAULTS_SYM                        = 436
	m57_FETCH_SYM                         = 437
	m57_FILE_SYM                          = 438
	m57_FILE_BLOCK_SIZE_SYM               = 439
	m57_FILTER_SYM                        = 440
	m57_FIRST_SYM                         = 441
	m57_FIXED_SYM                         = 442
	m57_FLOAT_NUM                         = 443
	m57_FLOAT_SYM                         = 444
	m57_FLUSH_SYM                         = 445
	m57_FOLLOWS_SYM                       = 446
	m57_FORCE_SYM                         = 447
	m57_FOREIGN                           = 448
	m57_FOR_SYM                           = 449
	m57_FORMAT_SYM                        = 450
	m57_FOUND_SYM                         = 451
	m57_FROM                              = 452
	m57_FULL                              = 453
	m57_FULLTEXT_SYM                      = 454
	m57_FUNCTION_SYM                      = 455
	m57_GE                                = 456
	m57_GENERAL                           = 457
	m57_GENERATED                         = 458
	m57_GROUP_REPLICATION                 = 459
	m57_GEOMETRYCOLLECTION                = 460
	m57_GEOMETRY_SYM                      = 461
	m57_GET_FORMAT                        = 462
	m57_GET_SYM                           = 463
	m57_GLOBAL_SYM                        = 464
	m57_GRANT                             = 465
	m57_GRANTS                            = 466
	m57_GROUP_SYM                         = 467
	m57_GROUP_CONCAT_SYM                  = 468
	m57_GT_SYM                            = 469
	m57_HANDLER_SYM                       = 470
	m57_HASH_SYM                          = 471
	m57_HAVING                            = 472
	m57_HELP_SYM                          = 473
	m57_HEX_NUM                           = 474
	m57_HIGH_PRIORITY                     = 475
	m57_HOST_SYM                          = 476
	m57_HOSTS_SYM                         = 477
	m57_HOUR_MICROSECOND_SYM              = 478
	m57_HOUR_MINUTE_SYM                   = 479
	m57_HOUR_SECOND_SYM                   = 480
	m57_HOUR_SYM                          = 481
	m57_IDENT                             = 482
	m57_IDENTIFIED_SYM                    = 483
	m57_IDENT_QUOTED                      = 484
	m57_IF                                = 485
	m57_IGNORE_SYM                        = 486
	m57_IGNORE_SERVER_IDS_SYM             = 487
	m57_IMPORT                            = 488
	m57_INDEXES                           = 489
	m57_INDEX_SYM                         = 490
	m57_INFILE                            = 491
	m57_INITIAL_SIZE_SYM                  = 492
	m57_INNER_SYM                         = 493
	m57_INOUT_SYM                         = 494
	m57_INSENSITIVE_SYM                   = 495
	m57_INSERT                            = 496
	m57_INSERT_METHOD                     = 497
	m57_INSTANCE_SYM                      = 498
	m57_INSTALL_SYM                       = 499
	m57_INTERVAL_SYM                      = 500
	m57_INTO                              = 501
	m57_INT_SYM                           = 502
	m57_INVOKER_SYM                       = 503
	m57_IN_SYM                            = 504
	m57_IO_AFTER_GTIDS                    = 505
	m57_IO_BEFORE_GTIDS                   = 506
	m57_IO_SYM                            = 507
	m57_IPC_SYM                           = 508
	m57_IS                                = 509
	m57_ISOLATION                         = 510
	m57_ISSUER_SYM                        = 511
	m57_ITERATE_SYM                       = 512
	m57_JOIN_SYM                          = 513
	m57_JSON_SEPARATOR_SYM                = 514
	m57_JSON_UNQUOTED_SEPARATOR_SYM       = 515
	m57_JSON_SYM                          = 516
	m57_KEYS                              = 517
	m57_KEY_BLOCK_SIZE                    = 518
	m57_KEY_SYM                           = 519
	m57_KILL_SYM                          = 520
	m57_LANGUAGE_SYM                      = 521
	m57_LAST_SYM                          = 522
	m57_LE                                = 523
	m57_LEADING                           = 524
	m57_LEAVES                            = 525
	m57_LEAVE_SYM                         = 526
	m57_LEFT                              = 527
	m57_LESS_SYM                          = 528
	m57_LEVEL_SYM                         = 529
	m57_LEX_HOSTNAME                      = 530
	m57_LIKE                              = 531
	m57_LIMIT                             = 532
	m57_LINEAR_SYM                        = 533
	m57_LINES                             = 534
	m57_LINESTRING                        = 535
	m57_LIST_SYM                          = 536
	m57_LOAD                              = 537
	m57_LOCAL_SYM                         = 538
	m57_LOCATOR_SYM                       = 539
	m57_LOCKS_SYM                         = 540
	m57_LOCK_SYM                          = 541
	m57_LOGFILE_SYM                       = 542
	m57_LOGS_SYM                          = 543
	m57_LONGBLOB                          = 544
	m57_LONGTEXT                          = 545
	m57_LONG_NUM                          = 546
	m57_LONG_SYM                          = 547
	m57_LOOP_SYM                          = 548
	m57_LOW_PRIORITY                      = 549
	m57_LT                                = 550
	m57_MASTER_AUTO_POSITION_SYM          = 551
	m57_MASTER_BIND_SYM                   = 552
	m57_MASTER_CONNECT_RETRY_SYM          = 553
	m57_MASTER_DELAY_SYM                  = 554
	m57_MASTER_HOST_SYM                   = 555
	m57_MASTER_LOG_FILE_SYM               = 556
	m57_MASTER_LOG_POS_SYM                = 557
	m57_MASTER_PASSWORD_SYM               = 558
	m57_MASTER_PORT_SYM                   = 559
	m57_MASTER_RETRY_COUNT_SYM            = 560
	m57_MASTER_SERVER_ID_SYM              = 561
	m57_MASTER_SSL_CAPATH_SYM             = 562
	m57_MASTER_TLS_VERSION_SYM            = 563
	m57_MASTER_SSL_CA_SYM                 = 564
	m57_MASTER_SSL_CERT_SYM               = 565
	m57_MASTER_SSL_CIPHER_SYM             = 566
	m57_MASTER_SSL_CRL_SYM                = 567
	m57_MASTER_SSL_CRLPATH_SYM            = 568
	m57_MASTER_SSL_KEY_SYM                = 569
	m57_MASTER_SSL_SYM                    = 570
	m57_MASTER_SSL_VERIFY_SERVER_CERT_SYM = 571
	m57_MASTER_SYM                        = 572
	m57_MASTER_USER_SYM                   = 573
	m57_MASTER_HEARTBEAT_PERIOD_SYM       = 574
	m57_MATCH                             = 575
	m57_MAX_CONNECTIONS_PER_HOUR          = 576
	m57_MAX_QUERIES_PER_HOUR              = 577
	m57_MAX_ROWS                          = 578
	m57_MAX_SIZE_SYM                      = 579
	m57_MAX_SYM                           = 580
	m57_MAX_UPDATES_PER_HOUR              = 581
	m57_MAX_USER_CONNECTIONS_SYM          = 582
	m57_MAX_VALUE_SYM                     = 583
	m57_MEDIUMBLOB                        = 584
	m57_MEDIUMINT                         = 585
	m57_MEDIUMTEXT                        = 586
	m57_MEDIUM_SYM                        = 587
	m57_MEMORY_SYM                        = 588
	m57_MERGE_SYM                         = 589
	m57_MESSAGE_TEXT_SYM                  = 590
	m57_MICROSECOND_SYM                   = 591
	m57_MIGRATE_SYM                       = 592
	m57_MINUTE_MICROSECOND_SYM            = 593
	m57_MINUTE_SECOND_SYM                 = 594
	m57_MINUTE_SYM                        = 595
	m57_MIN_ROWS                          = 596
	m57_MIN_SYM                           = 597
	m57_MODE_SYM                          = 598
	m57_MODIFIES_SYM                      = 599
	m57_MODIFY_SYM                        = 600
	m57_MOD_SYM                           = 601
	m57_MONTH_SYM                         = 602
	m57_MULTILINESTRING                   = 603
	m57_MULTIPOINT                        = 604
	m57_MULTIPOLYGON                      = 605
	m57_MUTEX_SYM                         = 606
	m57_MYSQL_ERRNO_SYM                   = 607
	m57_NAMES_SYM                         = 608
	m57_NAME_SYM                          = 609
	m57_NATIONAL_SYM                      = 610
	m57_NATURAL                           = 611
	m57_NCHAR_STRING                      = 612
	m57_NCHAR_SYM                         = 613
	m57_NDBCLUSTER_SYM                    = 614
	m57_NE                                = 615
	m57_NEG                               = 616
	m57_NEVER_SYM                         = 617
	m57_NEW_SYM                           = 618
	m57_NEXT_SYM                          = 619
	m57_NODEGROUP_SYM                     = 620
	m57_NONE_SYM                          = 621
	m57_NOT2_SYM                          = 622
	m57_NOT_SYM                           = 623
	m57_NOW_SYM                           = 624
	m57_NO_SYM                            = 625
	m57_NO_WAIT_SYM                       = 626
	m57_NO_WRITE_TO_BINLOG                = 627
	m57_NULL_SYM                          = 628
	m57_NUM                               = 629
	m57_NUMBER_SYM                        = 630
	m57_NUMERIC_SYM                       = 631
	m57_NVARCHAR_SYM                      = 632
	m57_OFFSET_SYM                        = 633
	m57_ON                                = 634
	m57_ONE_SYM                           = 635
	m57_ONLY_SYM                          = 636
	m57_OPEN_SYM                          = 637
	m57_OPTIMIZE                          = 638
	m57_OPTIMIZER_COSTS_SYM               = 639
	m57_OPTIONS_SYM                       = 640
	m57_OPTION                            = 641
	m57_OPTIONALLY                        = 642
	m57_OR2_SYM                           = 643
	m57_ORDER_SYM                         = 644
	m57_OR_OR_SYM                         = 645
	m57_OR_SYM                            = 646
	m57_OUTER                             = 647
	m57_OUTFILE                           = 648
	m57_OUT_SYM                           = 649
	m57_OWNER_SYM                         = 650
	m57_PACK_KEYS_SYM                     = 651
	m57_PAGE_SYM                          = 652
	m57_PARAM_MARKER                      = 653
	m57_PARSER_SYM                        = 654
	m57_PARSE_GCOL_EXPR_SYM               = 655
	m57_PARTIAL                           = 656
	m57_PARTITION_SYM                     = 657
	m57_PARTITIONS_SYM                    = 658
	m57_PARTITIONING_SYM                  = 659
	m57_PASSWORD                          = 660
	m57_PHASE_SYM                         = 661
	m57_PLUGIN_DIR_SYM                    = 662
	m57_PLUGIN_SYM                        = 663
	m57_PLUGINS_SYM                       = 664
	m57_POINT_SYM                         = 665
	m57_POLYGON                           = 666
	m57_PORT_SYM                          = 667
	m57_POSITION_SYM                      = 668
	m57_PRECEDES_SYM                      = 669
	m57_PRECISION                         = 670
	m57_PREPARE_SYM                       = 671
	m57_PRESERVE_SYM                      = 672
	m57_PREV_SYM                          = 673
	m57_PRIMARY_SYM                       = 674
	m57_PRIVILEGES                        = 675
	m57_PROCEDURE_SYM                     = 676
	m57_PROCESS                           = 677
	m57_PROCESSLIST_SYM                   = 678
	m57_PROFILE_SYM                       = 679
	m57_PROFILES_SYM                      = 680
	m57_PROXY_SYM                         = 681
	m57_PURGE                             = 682
	m57_QUARTER_SYM                       = 683
	m57_QUERY_SYM                         = 684
	m57_QUICK                             = 685
	m57_RANGE_SYM                         = 686
	m57_READS_SYM                         = 687
	m57_READ_ONLY_SYM                     = 688
	m57_READ_SYM                          = 689
	m57_READ_WRITE_SYM                    = 690
	m57_REAL                              = 691
	m57_REBUILD_SYM                       = 692
	m57_RECOVER_SYM                       = 693
	m57_REDOFILE_SYM                      = 694
	m57_REDO_BUFFER_SIZE_SYM              = 695
	m57_REDUNDANT_SYM                     = 696
	m57_REFERENCES                        = 697
	m57_REGEXP                            = 698
	m57_RELAY                             = 699
	m57_RELAYLOG_SYM                      = 700
	m57_RELAY_LOG_FILE_SYM                = 701
	m57_RELAY_LOG_POS_SYM                 = 702
	m57_RELAY_THREAD                      = 703
	m57_RELEASE_SYM                       = 704
	m57_RELOAD                            = 705
	m57_REMOVE_SYM                        = 706
	m57_RENAME                            = 707
	m57_REORGANIZE_SYM                    = 708
	m57_REPAIR                            = 709
	m57_REPEATABLE_SYM                    = 710
	m57_REPEAT_SYM                        = 711
	m57_REPLACE                           = 712
	m57_REPLICATION                       = 713
	m57_REPLICATE_DO_DB                   = 714
	m57_REPLICATE_IGNORE_DB               = 715
	m57_REPLICATE_DO_TABLE                = 716
	m57_REPLICATE_IGNORE_TABLE            = 717
	m57_REPLICATE_WILD_DO_TABLE           = 718
	m57_REPLICATE_WILD_IGNORE_TABLE       = 719
	m57_REPLICATE_REWRITE_DB              = 720
	m57_REQUIRE_SYM                       = 721
	m57_RESET_SYM                         = 722
	m57_RESIGNAL_SYM                      = 723
	m57_RESOURCES                         = 724
	m57_RESTORE_SYM                       = 725
	m57_RESTRICT                          = 726
	m57_RESUME_SYM                        = 727
	m57_RETURNED_SQLSTATE_SYM             = 728
	m57_RETURNS_SYM                       = 729
	m57_RETURN_SYM                        = 730
	m57_REVERSE_SYM                       = 731
	m57_REVOKE                            = 732
	m57_RIGHT                             = 733
	m57_ROLLBACK_SYM                      = 734
	m57_ROLLUP_SYM                        = 735
	m57_ROTATE_SYM                        = 736
	m57_ROUTINE_SYM                       = 737
	m57_ROWS_SYM                          = 738
	m57_ROW_FORMAT_SYM                    = 739
	m57_ROW_SYM                           = 740
	m57_ROW_COUNT_SYM                     = 741
	m57_RTREE_SYM                         = 742
	m57_SAVEPOINT_SYM                     = 743
	m57_SCHEDULE_SYM                      = 744
	m57_SCHEMA_NAME_SYM                   = 745
	m57_SECOND_MICROSECOND_SYM            = 746
	m57_SECOND_SYM                        = 747
	m57_SECURITY_SYM                      = 748
	m57_SELECT_SYM                        = 749
	m57_SENSITIVE_SYM                     = 750
	m57_SEPARATOR_SYM                     = 751
	m57_SERIALIZABLE_SYM                  = 752
	m57_SERIAL_SYM                        = 753
	m57_SESSION_SYM                       = 754
	m57_SERVER_SYM                        = 755
	m57_SERVER_OPTIONS                    = 756
	m57_SET                               = 757
	m57_SET_VAR                           = 758
	m57_SHARE_SYM                         = 759
	m57_SHIFT_LEFT                        = 760
	m57_SHIFT_RIGHT                       = 761
	m57_SHOW                              = 762
	m57_SHUTDOWN                          = 763
	m57_SIGNAL_SYM                        = 764
	m57_SIGNED_SYM                        = 765
	m57_SIMPLE_SYM                        = 766
	m57_SLAVE                             = 767
	m57_SLOW                              = 768
	m57_SMALLINT                          = 769
	m57_SNAPSHOT_SYM                      = 770
	m57_SOCKET_SYM                        = 771
	m57_SONAME_SYM                        = 772
	m57_SOUNDS_SYM                        = 773
	m57_SOURCE_SYM                        = 774
	m57_SPATIAL_SYM                       = 775
	m57_SPECIFIC_SYM                      = 776
	m57_SQLEXCEPTION_SYM                  = 777
	m57_SQLSTATE_SYM                      = 778
	m57_SQLWARNING_SYM                    = 779
	m57_SQL_AFTER_GTIDS                   = 780
	m57_SQL_AFTER_MTS_GAPS                = 781
	m57_SQL_BEFORE_GTIDS                  = 782
	m57_SQL_BIG_RESULT                    = 783
	m57_SQL_BUFFER_RESULT                 = 784
	m57_SQL_CACHE_SYM                     = 785
	m57_SQL_CALC_FOUND_ROWS               = 786
	m57_SQL_NO_CACHE_SYM                  = 787
	m57_SQL_SMALL_RESULT                  = 788
	m57_SQL_SYM                           = 789
	m57_SQL_THREAD                        = 790
	m57_SSL_SYM                           = 791
	m57_STACKED_SYM                       = 792
	m57_STARTING                          = 793
	m57_STARTS_SYM                        = 794
	m57_START_SYM                         = 795
	m57_STATS_AUTO_RECALC_SYM             = 796
	m57_STATS_PERSISTENT_SYM              = 797
	m57_STATS_SAMPLE_PAGES_SYM            = 798
	m57_STATUS_SYM                        = 799
	m57_STDDEV_SAMP_SYM                   = 800
	m57_STD_SYM                           = 801
	m57_STOP_SYM                          = 802
	m57_STORAGE_SYM                       = 803
	m57_STORED_SYM                        = 804
	m57_STRAIGHT_JOIN                     = 805
	m57_STRING_SYM                        = 806
	m57_SUBCLASS_ORIGIN_SYM               = 807
	m57_SUBDATE_SYM                       = 808
	m57_SUBJECT_SYM                       = 809
	m57_SUBPARTITIONS_SYM                 = 810
	m57_SUBPARTITION_SYM                  = 811
	m57_SUBSTRING                         = 812
	m57_SUM_SYM                           = 813
	m57_SUPER_SYM                         = 814
	m57_SUSPEND_SYM                       = 815
	m57_SWAPS_SYM                         = 816
	m57_SWITCHES_SYM                      = 817
	m57_SYSDATE                           = 818
	m57_TABLES                            = 819
	m57_TABLESPACE_SYM                    = 820
	m57_TABLE_REF_PRIORITY                = 821
	m57_TABLE_SYM                         = 822
	m57_TABLE_CHECKSUM_SYM                = 823
	m57_TABLE_NAME_SYM                    = 824
	m57_TEMPORARY                         = 825
	m57_TEMPTABLE_SYM                     = 826
	m57_TERMINATED                        = 827
	m57_TEXT_STRING                       = 828
	m57_TEXT_SYM                          = 829
	m57_THAN_SYM                          = 830
	m57_THEN_SYM                          = 831
	m57_TIMESTAMP                         = 832
	m57_TIMESTAMP_ADD                     = 833
	m57_TIMESTAMP_DIFF                    = 834
	m57_TIME_SYM                          = 835
	m57_TINYBLOB                          = 836
	m57_TINYINT                           = 837
	m57_TINYTEXT                          = 838
	m57_TO_SYM                            = 839
	m57_TRAILING                          = 840
	m57_TRANSACTION_SYM                   = 841
	m57_TRIGGERS_SYM                      = 842
	m57_TRIGGER_SYM                       = 843
	m57_TRIM                              = 844
	m57_TRUE_SYM                          = 845
	m57_TRUNCATE_SYM                      = 846
	m57_TYPES_SYM                         = 847
	m57_TYPE_SYM                          = 848
	m57_UDF_RETURNS_SYM                   = 849
	m57_ULONGLONG_NUM                     = 850
	m57_UNCOMMITTED_SYM                   = 851
	m57_UNDEFINED_SYM                     = 852
	m57_UNDERSCORE_CHARSET                = 853
	m57_UNDOFILE_SYM                      = 854
	m57_UNDO_BUFFER_SIZE_SYM              = 855
	m57_UNDO_SYM                          = 856
	m57_UNICODE_SYM                       = 857
	m57_UNINSTALL_SYM                     = 858
	m57_UNION_SYM                         = 859
	m57_UNIQUE_SYM                        = 860
	m57_UNKNOWN_SYM                       = 861
	m57_UNLOCK_SYM                        = 862
	m57_UNSIGNED                          = 863
	m57_UNTIL_SYM                         = 864
	m57_UPDATE_SYM                        = 865
	m57_UPGRADE_SYM                       = 866
	m57_USAGE                             = 867
	m57_USER                              = 868
	m57_USE_FRM                           = 869
	m57_USE_SYM                           = 870
	m57_USING                             = 871
	m57_UTC_DATE_SYM                      = 872
	m57_UTC_TIMESTAMP_SYM                 = 873
	m57_UTC_TIME_SYM                      = 874
	m57_VALIDATION_SYM                    = 875
	m57_VALUES                            = 876
	m57_VALUE_SYM                         = 877
	m57_VARBINARY                         = 878
	m57_VARCHAR                           = 879
	m57_VARIABLES                         = 880
	m57_VARIANCE_SYM                      = 881
	m57_VARYING                           = 882
	m57_VAR_SAMP_SYM                      = 883
	m57_VIEW_SYM                          = 884
	m57_VIRTUAL_SYM                       = 885
	m57_WAIT_SYM                          = 886
	m57_WARNINGS                          = 887
	m57_WEEK_SYM                          = 888
	m57_WEIGHT_STRING_SYM                 = 889
	m57_WHEN_SYM                          = 890
	m57_WHERE                             = 891
	m57_WHILE_SYM                         = 892
	m57_WITH                              = 893
	m57_WITH_CUBE_SYM                     = 894
	m57_WITH_ROLLUP_SYM                   = 895
	m57_WITHOUT_SYM                       = 896
	m57_WORK_SYM                          = 897
	m57_WRAPPER_SYM                       = 898
	m57_WRITE_SYM                         = 899
	m57_X509_SYM                          = 900
	m57_XA_SYM                            = 901
	m57_XID_SYM                           = 902
	m57_XML_SYM                           = 903
	m57_XOR                               = 904
	m57_YEAR_MONTH_SYM                    = 905
	m57_YEAR_SYM                          = 906
	m57_ZEROFILL                          = 907
	m57_JSON_OBJECTAGG                    = 908
	m57_JSON_ARRAYAGG                     = 909

	// MySQL 5.7 hint tokens
	m57_MAX_EXECUTION_TIME_HINT    = 910
	m57_BKA_HINT                   = 911
	m57_BNL_HINT                   = 912
	m57_DUPSWEEDOUT_HINT           = 913
	m57_FIRSTMATCH_HINT            = 914
	m57_INTOEXISTS_HINT            = 915
	m57_LOOSESCAN_HINT             = 916
	m57_MATERIALIZATION_HINT       = 917
	m57_NO_BKA_HINT                = 918
	m57_NO_BNL_HINT                = 919
	m57_NO_ICP_HINT                = 920
	m57_NO_MRR_HINT                = 921
	m57_NO_RANGE_OPTIMIZATION_HINT = 922
	m57_NO_SEMIJOIN_HINT           = 923
	m57_MRR_HINT                   = 924
	m57_QB_NAME_HINT               = 925
	m57_SEMIJOIN_HINT              = 926
	m57_SUBQUERY_HINT              = 927
	m57_HINT_ARG_NUMBER            = 928
	m57_HINT_ARG_IDENT             = 929
	m57_HINT_ARG_QB_NAME           = 930
	m57_HINT_CLOSE                 = 931
	m57_HINT_ERROR                 = 932

	// MySQL 5.7 digest tokens
	m57TOK_GENERIC_VALUE           = 928
	m57TOK_GENERIC_VALUE_LIST      = 929
	m57TOK_ROW_SINGLE_VALUE        = 930
	m57TOK_ROW_SINGLE_VALUE_LIST   = 931
	m57TOK_ROW_MULTIPLE_VALUE      = 932
	m57TOK_ROW_MULTIPLE_VALUE_LIST = 933
	m57TOK_IDENT                   = 934
	m57TOK_IDENT_AT                = 935
	m57TOK_HINT_COMMENT_OPEN       = 936
	m57TOK_HINT_COMMENT_CLOSE      = 937
	m57TOK_UNUSED                  = 938
)

// mysql80To57TokenMap maps MySQL 8.0 token IDs to MySQL 5.7 token IDs.
var mysql80To57TokenMap = map[int]int{
	// Main SQL tokens (258+)
	// These are listed in MySQL 8.0 token order with their MySQL 5.7 mappings
	ABORT_SYM:                           m57_ABORT_SYM,
	ACCESSIBLE_SYM:                      m57_ACCESSIBLE_SYM,
	ACCOUNT_SYM:                         m57_ACCOUNT_SYM,
	ACTION:                              m57_ACTION,
	ADD:                                 m57_ADD,
	ADDDATE_SYM:                         m57_ADDDATE_SYM,
	AFTER_SYM:                           m57_AFTER_SYM,
	AGAINST:                             m57_AGAINST,
	AGGREGATE_SYM:                       m57_AGGREGATE_SYM,
	ALGORITHM_SYM:                       m57_ALGORITHM_SYM,
	ALL:                                 m57_ALL,
	ALTER:                               m57_ALTER,
	ALWAYS_SYM:                          m57_ALWAYS_SYM,
	OBSOLETE_TOKEN_271:                  m57_ANALYSE_SYM, // ANALYSE_SYM removed in 8.0
	ANALYZE_SYM:                         m57_ANALYZE_SYM,
	AND_AND_SYM:                         m57_AND_AND_SYM,
	AND_SYM:                             m57_AND_SYM,
	ANY_SYM:                             m57_ANY_SYM,
	AS:                                  m57_AS,
	ASC:                                 m57_ASC,
	ASCII_SYM:                           m57_ASCII_SYM,
	ASENSITIVE_SYM:                      m57_ASENSITIVE_SYM,
	AT_SYM:                              m57_AT_SYM,
	AUTOEXTEND_SIZE_SYM:                 m57_AUTOEXTEND_SIZE_SYM,
	AUTO_INC:                            m57_AUTO_INC,
	AVG_ROW_LENGTH:                      m57_AVG_ROW_LENGTH,
	AVG_SYM:                             m57_AVG_SYM,
	BACKUP_SYM:                          m57_BACKUP_SYM,
	BEFORE_SYM:                          m57_BEFORE_SYM,
	BEGIN_SYM:                           m57_BEGIN_SYM,
	BETWEEN_SYM:                         m57_BETWEEN_SYM,
	BIGINT_SYM:                          m57_BIGINT,
	BINARY_SYM:                          m57_BINARY,
	BINLOG_SYM:                          m57_BINLOG_SYM,
	BIN_NUM:                             m57_BIN_NUM,
	BIT_AND_SYM:                         m57_BIT_AND,
	BIT_OR_SYM:                          m57_BIT_OR,
	BIT_SYM:                             m57_BIT_SYM,
	BIT_XOR_SYM:                         m57_BIT_XOR,
	BLOB_SYM:                            m57_BLOB_SYM,
	BLOCK_SYM:                           m57_BLOCK_SYM,
	BOOLEAN_SYM:                         m57_BOOLEAN_SYM,
	BOOL_SYM:                            m57_BOOL_SYM,
	BOTH:                                m57_BOTH,
	BTREE_SYM:                           m57_BTREE_SYM,
	BY:                                  m57_BY,
	BYTE_SYM:                            m57_BYTE_SYM,
	CACHE_SYM:                           m57_CACHE_SYM,
	CALL_SYM:                            m57_CALL_SYM,
	CASCADE:                             m57_CASCADE,
	CASCADED:                            m57_CASCADED,
	CASE_SYM:                            m57_CASE_SYM,
	CAST_SYM:                            m57_CAST_SYM,
	CATALOG_NAME_SYM:                    m57_CATALOG_NAME_SYM,
	CHAIN_SYM:                           m57_CHAIN_SYM,
	CHANGE:                              m57_CHANGE,
	CHANGED:                             m57_CHANGED,
	CHANNEL_SYM:                         m57_CHANNEL_SYM,
	CHARSET:                             m57_CHARSET,
	CHAR_SYM:                            m57_CHAR_SYM,
	CHECKSUM_SYM:                        m57_CHECKSUM_SYM,
	CHECK_SYM:                           m57_CHECK_SYM,
	CIPHER_SYM:                          m57_CIPHER_SYM,
	CLASS_ORIGIN_SYM:                    m57_CLASS_ORIGIN_SYM,
	CLIENT_SYM:                          m57_CLIENT_SYM,
	CLOSE_SYM:                           m57_CLOSE_SYM,
	COALESCE:                            m57_COALESCE,
	CODE_SYM:                            m57_CODE_SYM,
	COLLATE_SYM:                         m57_COLLATE_SYM,
	COLLATION_SYM:                       m57_COLLATION_SYM,
	COLUMNS:                             m57_COLUMNS,
	COLUMN_SYM:                          m57_COLUMN_SYM,
	COLUMN_FORMAT_SYM:                   m57_COLUMN_FORMAT_SYM,
	COLUMN_NAME_SYM:                     m57_COLUMN_NAME_SYM,
	COMMENT_SYM:                         m57_COMMENT_SYM,
	COMMITTED_SYM:                       m57_COMMITTED_SYM,
	COMMIT_SYM:                          m57_COMMIT_SYM,
	COMPACT_SYM:                         m57_COMPACT_SYM,
	COMPLETION_SYM:                      m57_COMPLETION_SYM,
	COMPRESSED_SYM:                      m57_COMPRESSED_SYM,
	COMPRESSION_SYM:                     m57_COMPRESSION_SYM,
	ENCRYPTION_SYM:                      m57_ENCRYPTION_SYM,
	CONCURRENT:                          m57_CONCURRENT,
	CONDITION_SYM:                       m57_CONDITION_SYM,
	CONNECTION_SYM:                      m57_CONNECTION_SYM,
	CONSISTENT_SYM:                      m57_CONSISTENT_SYM,
	CONSTRAINT:                          m57_CONSTRAINT,
	CONSTRAINT_CATALOG_SYM:              m57_CONSTRAINT_CATALOG_SYM,
	CONSTRAINT_NAME_SYM:                 m57_CONSTRAINT_NAME_SYM,
	CONSTRAINT_SCHEMA_SYM:               m57_CONSTRAINT_SCHEMA_SYM,
	CONTAINS_SYM:                        m57_CONTAINS_SYM,
	CONTEXT_SYM:                         m57_CONTEXT_SYM,
	CONTINUE_SYM:                        m57_CONTINUE_SYM,
	CONVERT_SYM:                         m57_CONVERT_SYM,
	COUNT_SYM:                           m57_COUNT_SYM,
	CPU_SYM:                             m57_CPU_SYM,
	CREATE:                              m57_CREATE,
	CROSS:                               m57_CROSS,
	CUBE_SYM:                            m57_CUBE_SYM,
	CURDATE:                             m57_CURDATE,
	CURRENT_SYM:                         m57_CURRENT_SYM,
	CURRENT_USER:                        m57_CURRENT_USER,
	CURSOR_SYM:                          m57_CURSOR_SYM,
	CURSOR_NAME_SYM:                     m57_CURSOR_NAME_SYM,
	CURTIME:                             m57_CURTIME,
	DATABASE:                            m57_DATABASE,
	DATABASES:                           m57_DATABASES,
	DATAFILE_SYM:                        m57_DATAFILE_SYM,
	DATA_SYM:                            m57_DATA_SYM,
	DATETIME_SYM:                        m57_DATETIME,
	DATE_ADD_INTERVAL:                   m57_DATE_ADD_INTERVAL,
	DATE_SUB_INTERVAL:                   m57_DATE_SUB_INTERVAL,
	DATE_SYM:                            m57_DATE_SYM,
	DAY_HOUR_SYM:                        m57_DAY_HOUR_SYM,
	DAY_MICROSECOND_SYM:                 m57_DAY_MICROSECOND_SYM,
	DAY_MINUTE_SYM:                      m57_DAY_MINUTE_SYM,
	DAY_SECOND_SYM:                      m57_DAY_SECOND_SYM,
	DAY_SYM:                             m57_DAY_SYM,
	DEALLOCATE_SYM:                      m57_DEALLOCATE_SYM,
	DECIMAL_NUM:                         m57_DECIMAL_NUM,
	DECIMAL_SYM:                         m57_DECIMAL_SYM,
	DECLARE_SYM:                         m57_DECLARE_SYM,
	DEFAULT_SYM:                         m57_DEFAULT,
	DEFAULT_AUTH_SYM:                    m57_DEFAULT_AUTH_SYM,
	DEFINER_SYM:                         m57_DEFINER_SYM,
	DELAYED_SYM:                         m57_DELAYED_SYM,
	DELAY_KEY_WRITE_SYM:                 m57_DELAY_KEY_WRITE_SYM,
	DELETE_SYM:                          m57_DELETE_SYM,
	DESC:                                m57_DESC,
	DESCRIBE:                            m57_DESCRIBE,
	OBSOLETE_TOKEN_388:                  m57_DES_KEY_FILE,
	DETERMINISTIC_SYM:                   m57_DETERMINISTIC_SYM,
	DIAGNOSTICS_SYM:                     m57_DIAGNOSTICS_SYM,
	DIRECTORY_SYM:                       m57_DIRECTORY_SYM,
	DISABLE_SYM:                         m57_DISABLE_SYM,
	DISCARD_SYM:                         m57_DISCARD,
	DISK_SYM:                            m57_DISK_SYM,
	DISTINCT:                            m57_DISTINCT,
	DIV_SYM:                             m57_DIV_SYM,
	DOUBLE_SYM:                          m57_DOUBLE_SYM,
	DO_SYM:                              m57_DO_SYM,
	DROP:                                m57_DROP,
	DUAL_SYM:                            m57_DUAL_SYM,
	DUMPFILE:                            m57_DUMPFILE,
	DUPLICATE_SYM:                       m57_DUPLICATE_SYM,
	DYNAMIC_SYM:                         m57_DYNAMIC_SYM,
	EACH_SYM:                            m57_EACH_SYM,
	ELSE:                                m57_ELSE,
	ELSEIF_SYM:                          m57_ELSEIF_SYM,
	ENABLE_SYM:                          m57_ENABLE_SYM,
	ENCLOSED:                            m57_ENCLOSED,
	END:                                 m57_END,
	ENDS_SYM:                            m57_ENDS_SYM,
	END_OF_INPUT:                        m57_END_OF_INPUT,
	ENGINES_SYM:                         m57_ENGINES_SYM,
	ENGINE_SYM:                          m57_ENGINE_SYM,
	ENUM_SYM:                            m57_ENUM,
	EQ:                                  m57_EQ,
	EQUAL_SYM:                           m57_EQUAL_SYM,
	ERROR_SYM:                           m57_ERROR_SYM,
	ERRORS:                              m57_ERRORS,
	ESCAPED:                             m57_ESCAPED,
	ESCAPE_SYM:                          m57_ESCAPE_SYM,
	EVENTS_SYM:                          m57_EVENTS_SYM,
	EVENT_SYM:                           m57_EVENT_SYM,
	EVERY_SYM:                           m57_EVERY_SYM,
	EXCHANGE_SYM:                        m57_EXCHANGE_SYM,
	EXECUTE_SYM:                         m57_EXECUTE_SYM,
	EXISTS:                              m57_EXISTS,
	EXIT_SYM:                            m57_EXIT_SYM,
	EXPANSION_SYM:                       m57_EXPANSION_SYM,
	EXPIRE_SYM:                          m57_EXPIRE_SYM,
	EXPORT_SYM:                          m57_EXPORT_SYM,
	EXTENDED_SYM:                        m57_EXTENDED_SYM,
	EXTENT_SIZE_SYM:                     m57_EXTENT_SIZE_SYM,
	EXTRACT_SYM:                         m57_EXTRACT_SYM,
	FALSE_SYM:                           m57_FALSE_SYM,
	FAST_SYM:                            m57_FAST_SYM,
	FAULTS_SYM:                          m57_FAULTS_SYM,
	FETCH_SYM:                           m57_FETCH_SYM,
	FILE_SYM:                            m57_FILE_SYM,
	FILE_BLOCK_SIZE_SYM:                 m57_FILE_BLOCK_SIZE_SYM,
	FILTER_SYM:                          m57_FILTER_SYM,
	FIRST_SYM:                           m57_FIRST_SYM,
	FIXED_SYM:                           m57_FIXED_SYM,
	FLOAT_NUM:                           m57_FLOAT_NUM,
	FLOAT_SYM:                           m57_FLOAT_SYM,
	FLUSH_SYM:                           m57_FLUSH_SYM,
	FOLLOWS_SYM:                         m57_FOLLOWS_SYM,
	FORCE_SYM:                           m57_FORCE_SYM,
	FOREIGN:                             m57_FOREIGN,
	FOR_SYM:                             m57_FOR_SYM,
	FORMAT_SYM:                          m57_FORMAT_SYM,
	FOUND_SYM:                           m57_FOUND_SYM,
	FROM:                                m57_FROM,
	FULL:                                m57_FULL,
	FULLTEXT_SYM:                        m57_FULLTEXT_SYM,
	FUNCTION_SYM:                        m57_FUNCTION_SYM,
	GE:                                  m57_GE,
	GENERAL:                             m57_GENERAL,
	GENERATED:                           m57_GENERATED,
	GROUP_REPLICATION:                   m57_GROUP_REPLICATION,
	GEOMETRYCOLLECTION_SYM:              m57_GEOMETRYCOLLECTION,
	GEOMETRY_SYM:                        m57_GEOMETRY_SYM,
	GET_FORMAT:                          m57_GET_FORMAT,
	GET_SYM:                             m57_GET_SYM,
	GLOBAL_SYM:                          m57_GLOBAL_SYM,
	GRANT:                               m57_GRANT,
	GRANTS:                              m57_GRANTS,
	GROUP_SYM:                           m57_GROUP_SYM,
	GROUP_CONCAT_SYM:                    m57_GROUP_CONCAT_SYM,
	GT_SYM:                              m57_GT_SYM,
	HANDLER_SYM:                         m57_HANDLER_SYM,
	HASH_SYM:                            m57_HASH_SYM,
	HAVING:                              m57_HAVING,
	HELP_SYM:                            m57_HELP_SYM,
	HEX_NUM:                             m57_HEX_NUM,
	HIGH_PRIORITY:                       m57_HIGH_PRIORITY,
	HOST_SYM:                            m57_HOST_SYM,
	HOSTS_SYM:                           m57_HOSTS_SYM,
	HOUR_MICROSECOND_SYM:                m57_HOUR_MICROSECOND_SYM,
	HOUR_MINUTE_SYM:                     m57_HOUR_MINUTE_SYM,
	HOUR_SECOND_SYM:                     m57_HOUR_SECOND_SYM,
	HOUR_SYM:                            m57_HOUR_SYM,
	IDENT:                               m57_IDENT,
	IDENTIFIED_SYM:                      m57_IDENTIFIED_SYM,
	IDENT_QUOTED:                        m57_IDENT_QUOTED,
	IF:                                  m57_IF,
	IGNORE_SYM:                          m57_IGNORE_SYM,
	IGNORE_SERVER_IDS_SYM:               m57_IGNORE_SERVER_IDS_SYM,
	IMPORT:                              m57_IMPORT,
	INDEXES:                             m57_INDEXES,
	INDEX_SYM:                           m57_INDEX_SYM,
	INFILE_SYM:                          m57_INFILE,
	INITIAL_SIZE_SYM:                    m57_INITIAL_SIZE_SYM,
	INNER_SYM:                           m57_INNER_SYM,
	INOUT_SYM:                           m57_INOUT_SYM,
	INSENSITIVE_SYM:                     m57_INSENSITIVE_SYM,
	INSERT_SYM:                          m57_INSERT,
	INSERT_METHOD:                       m57_INSERT_METHOD,
	INSTANCE_SYM:                        m57_INSTANCE_SYM,
	INSTALL_SYM:                         m57_INSTALL_SYM,
	INTERVAL_SYM:                        m57_INTERVAL_SYM,
	INTO:                                m57_INTO,
	INT_SYM:                             m57_INT_SYM,
	INVOKER_SYM:                         m57_INVOKER_SYM,
	IN_SYM:                              m57_IN_SYM,
	IO_AFTER_GTIDS:                      m57_IO_AFTER_GTIDS,
	IO_BEFORE_GTIDS:                     m57_IO_BEFORE_GTIDS,
	IO_SYM:                              m57_IO_SYM,
	IPC_SYM:                             m57_IPC_SYM,
	IS:                                  m57_IS,
	ISOLATION:                           m57_ISOLATION,
	ISSUER_SYM:                          m57_ISSUER_SYM,
	ITERATE_SYM:                         m57_ITERATE_SYM,
	JOIN_SYM:                            m57_JOIN_SYM,
	JSON_SEPARATOR_SYM:                  m57_JSON_SEPARATOR_SYM,
	JSON_SYM:                            m57_JSON_SYM,
	KEYS:                                m57_KEYS,
	KEY_BLOCK_SIZE:                      m57_KEY_BLOCK_SIZE,
	KEY_SYM:                             m57_KEY_SYM,
	KILL_SYM:                            m57_KILL_SYM,
	LANGUAGE_SYM:                        m57_LANGUAGE_SYM,
	LAST_SYM:                            m57_LAST_SYM,
	LE:                                  m57_LE,
	LEADING:                             m57_LEADING,
	LEAVES:                              m57_LEAVES,
	LEAVE_SYM:                           m57_LEAVE_SYM,
	LEFT:                                m57_LEFT,
	LESS_SYM:                            m57_LESS_SYM,
	LEVEL_SYM:                           m57_LEVEL_SYM,
	LEX_HOSTNAME:                        m57_LEX_HOSTNAME,
	LIKE:                                m57_LIKE,
	LIMIT:                               m57_LIMIT,
	LINEAR_SYM:                          m57_LINEAR_SYM,
	LINES:                               m57_LINES,
	LINESTRING_SYM:                      m57_LINESTRING,
	LIST_SYM:                            m57_LIST_SYM,
	LOAD:                                m57_LOAD,
	LOCAL_SYM:                           m57_LOCAL_SYM,
	OBSOLETE_TOKEN_538:                  m57_LOCATOR_SYM,
	LOCKS_SYM:                           m57_LOCKS_SYM,
	LOCK_SYM:                            m57_LOCK_SYM,
	LOGFILE_SYM:                         m57_LOGFILE_SYM,
	LOGS_SYM:                            m57_LOGS_SYM,
	LONGBLOB_SYM:                        m57_LONGBLOB,
	LONGTEXT_SYM:                        m57_LONGTEXT,
	LONG_NUM:                            m57_LONG_NUM,
	LONG_SYM:                            m57_LONG_SYM,
	LOOP_SYM:                            m57_LOOP_SYM,
	LOW_PRIORITY:                        m57_LOW_PRIORITY,
	LT:                                  m57_LT,
	OBSOLETE_TOKEN_550:                  m57_MASTER_AUTO_POSITION_SYM,
	OBSOLETE_TOKEN_551:                  m57_MASTER_BIND_SYM,
	OBSOLETE_TOKEN_552:                  m57_MASTER_CONNECT_RETRY_SYM,
	OBSOLETE_TOKEN_553:                  m57_MASTER_DELAY_SYM,
	OBSOLETE_TOKEN_554:                  m57_MASTER_HOST_SYM,
	OBSOLETE_TOKEN_555:                  m57_MASTER_LOG_FILE_SYM,
	OBSOLETE_TOKEN_556:                  m57_MASTER_LOG_POS_SYM,
	OBSOLETE_TOKEN_557:                  m57_MASTER_PASSWORD_SYM,
	OBSOLETE_TOKEN_558:                  m57_MASTER_PORT_SYM,
	OBSOLETE_TOKEN_559:                  m57_MASTER_RETRY_COUNT_SYM,
	OBSOLETE_TOKEN_561:                  m57_MASTER_SERVER_ID_SYM,
	OBSOLETE_TOKEN_562:                  m57_MASTER_SSL_SYM,
	OBSOLETE_TOKEN_563:                  m57_MASTER_SSL_CA_SYM,
	OBSOLETE_TOKEN_564:                  m57_MASTER_SSL_CAPATH_SYM,
	OBSOLETE_TOKEN_565:                  m57_MASTER_SSL_CERT_SYM,
	OBSOLETE_TOKEN_566:                  m57_MASTER_SSL_CIPHER_SYM,
	OBSOLETE_TOKEN_567:                  m57_MASTER_SSL_CRL_SYM,
	OBSOLETE_TOKEN_568:                  m57_MASTER_SSL_CRLPATH_SYM,
	OBSOLETE_TOKEN_569:                  m57_MASTER_SSL_KEY_SYM,
	OBSOLETE_TOKEN_570:                  m57_MASTER_SSL_VERIFY_SERVER_CERT_SYM,
	MASTER_SYM:                          m57_MASTER_SYM,
	OBSOLETE_TOKEN_572:                  m57_MASTER_TLS_VERSION_SYM,
	OBSOLETE_TOKEN_573:                  m57_MASTER_USER_SYM,
	MATCH:                               m57_MATCH,
	MAX_CONNECTIONS_PER_HOUR:            m57_MAX_CONNECTIONS_PER_HOUR,
	MAX_QUERIES_PER_HOUR:                m57_MAX_QUERIES_PER_HOUR,
	MAX_ROWS:                            m57_MAX_ROWS,
	MAX_SIZE_SYM:                        m57_MAX_SIZE_SYM,
	MAX_SYM:                             m57_MAX_SYM,
	MAX_UPDATES_PER_HOUR:                m57_MAX_UPDATES_PER_HOUR,
	MAX_USER_CONNECTIONS_SYM:            m57_MAX_USER_CONNECTIONS_SYM,
	MAX_VALUE_SYM:                       m57_MAX_VALUE_SYM,
	MEDIUMBLOB_SYM:                      m57TOK_UNUSED, // Not in MySQL 5.7
	MEDIUMINT_SYM:                       m57TOK_UNUSED, // Not in MySQL 5.7
	MEDIUMTEXT_SYM:                      m57TOK_UNUSED, // Not in MySQL 5.7
	MEDIUM_SYM:                          m57_MEDIUM_SYM,
	MEMORY_SYM:                          m57_MEMORY_SYM,
	MERGE_SYM:                           m57_MERGE_SYM,
	MESSAGE_TEXT_SYM:                    m57_MESSAGE_TEXT_SYM,
	MICROSECOND_SYM:                     m57_MICROSECOND_SYM,
	MIGRATE_SYM:                         m57_MIGRATE_SYM,
	MINUTE_MICROSECOND_SYM:              m57_MINUTE_MICROSECOND_SYM,
	MINUTE_SECOND_SYM:                   m57_MINUTE_SECOND_SYM,
	MINUTE_SYM:                          m57_MINUTE_SYM,
	MIN_ROWS:                            m57_MIN_ROWS,
	MIN_SYM:                             m57_MIN_SYM,
	MODE_SYM:                            m57_MODE_SYM,
	MODIFIES_SYM:                        m57_MODIFIES_SYM,
	MODIFY_SYM:                          m57_MODIFY_SYM,
	MOD_SYM:                             m57_MOD_SYM,
	MONTH_SYM:                           m57_MONTH_SYM,
	MULTILINESTRING_SYM:                 m57TOK_UNUSED, // Not in MySQL 5.7
	MULTIPOINT_SYM:                      m57TOK_UNUSED, // Not in MySQL 5.7
	MULTIPOLYGON_SYM:                    m57TOK_UNUSED, // Not in MySQL 5.7
	MUTEX_SYM:                           m57_MUTEX_SYM,
	MYSQL_ERRNO_SYM:                     m57_MYSQL_ERRNO_SYM,
	NAMES_SYM:                           m57_NAMES_SYM,
	NAME_SYM:                            m57_NAME_SYM,
	NATIONAL_SYM:                        m57_NATIONAL_SYM,
	NATURAL:                             m57_NATURAL,
	NCHAR_STRING:                        m57_NCHAR_STRING,
	NCHAR_SYM:                           m57_NCHAR_SYM,
	NDBCLUSTER_SYM:                      m57_NDBCLUSTER_SYM,
	NE:                                  m57_NE,
	NEG:                                 m57_NEG,
	NEVER_SYM:                           m57_NEVER_SYM,
	NEW_SYM:                             m57_NEW_SYM,
	NEXT_SYM:                            m57_NEXT_SYM,
	NODEGROUP_SYM:                       m57_NODEGROUP_SYM,
	NONE_SYM:                            m57_NONE_SYM,
	NOT2_SYM:                            m57_NOT2_SYM,
	NOT_SYM:                             m57_NOT_SYM,
	NOW_SYM:                             m57_NOW_SYM,
	NO_SYM:                              m57_NO_SYM,
	NO_WAIT_SYM:                         m57_NO_WAIT_SYM,
	NO_WRITE_TO_BINLOG:                  m57_NO_WRITE_TO_BINLOG,
	NULL_SYM:                            m57_NULL_SYM,
	NUM:                                 m57_NUM,
	NUMBER_SYM:                          m57_NUMBER_SYM,
	NUMERIC_SYM:                         m57_NUMERIC_SYM,
	NVARCHAR_SYM:                        m57_NVARCHAR_SYM,
	OFFSET_SYM:                          m57_OFFSET_SYM,
	ON_SYM:                              m57_ON,
	ONE_SYM:                             m57_ONE_SYM,
	ONLY_SYM:                            m57_ONLY_SYM,
	OPEN_SYM:                            m57_OPEN_SYM,
	OPTIMIZE:                            m57_OPTIMIZE,
	OPTIMIZER_COSTS_SYM:                 m57_OPTIMIZER_COSTS_SYM,
	OPTIONS_SYM:                         m57_OPTIONS_SYM,
	OPTION:                              m57_OPTION,
	OPTIONALLY:                          m57_OPTIONALLY,
	OR2_SYM:                             m57_OR2_SYM,
	ORDER_SYM:                           m57_ORDER_SYM,
	OR_OR_SYM:                           m57_OR_OR_SYM,
	OR_SYM:                              m57_OR_SYM,
	OUTER_SYM:                           m57_OUTER,
	OUTFILE:                             m57_OUTFILE,
	OUT_SYM:                             m57_OUT_SYM,
	OWNER_SYM:                           m57_OWNER_SYM,
	PACK_KEYS_SYM:                       m57_PACK_KEYS_SYM,
	PAGE_SYM:                            m57_PAGE_SYM,
	PARAM_MARKER:                        m57_PARAM_MARKER,
	PARSER_SYM:                          m57_PARSER_SYM,
	OBSOLETE_TOKEN_654:                  m57_PARSE_GCOL_EXPR_SYM,
	PARTIAL:                             m57_PARTIAL,
	PARTITION_SYM:                       m57_PARTITION_SYM,
	PARTITIONS_SYM:                      m57_PARTITIONS_SYM,
	PARTITIONING_SYM:                    m57_PARTITIONING_SYM,
	PASSWORD:                            m57_PASSWORD,
	PHASE_SYM:                           m57_PHASE_SYM,
	PLUGIN_DIR_SYM:                      m57_PLUGIN_DIR_SYM,
	PLUGIN_SYM:                          m57_PLUGIN_SYM,
	PLUGINS_SYM:                         m57_PLUGINS_SYM,
	POINT_SYM:                           m57_POINT_SYM,
	POLYGON_SYM:                         m57_POLYGON,
	PORT_SYM:                            m57_PORT_SYM,
	POSITION_SYM:                        m57_POSITION_SYM,
	PRECEDES_SYM:                        m57_PRECEDES_SYM,
	PRECISION:                           m57_PRECISION,
	PREPARE_SYM:                         m57_PREPARE_SYM,
	PRESERVE_SYM:                        m57_PRESERVE_SYM,
	PREV_SYM:                            m57_PREV_SYM,
	PRIMARY_SYM:                         m57_PRIMARY_SYM,
	PRIVILEGES:                          m57_PRIVILEGES,
	PROCEDURE_SYM:                       m57_PROCEDURE_SYM,
	PROCESS:                             m57_PROCESS,
	PROCESSLIST_SYM:                     m57_PROCESSLIST_SYM,
	PROFILE_SYM:                         m57_PROFILE_SYM,
	PROFILES_SYM:                        m57_PROFILES_SYM,
	PROXY_SYM:                           m57_PROXY_SYM,
	PURGE:                               m57_PURGE,
	QUARTER_SYM:                         m57_QUARTER_SYM,
	QUERY_SYM:                           m57_QUERY_SYM,
	QUICK:                               m57_QUICK,
	RANGE_SYM:                           m57_RANGE_SYM,
	READS_SYM:                           m57_READS_SYM,
	READ_ONLY_SYM:                       m57_READ_ONLY_SYM,
	READ_SYM:                            m57_READ_SYM,
	READ_WRITE_SYM:                      m57_READ_WRITE_SYM,
	REAL_SYM:                            m57_REAL,
	REBUILD_SYM:                         m57_REBUILD_SYM,
	RECOVER_SYM:                         m57_RECOVER_SYM,
	OBSOLETE_TOKEN_693:                  m57_REDOFILE_SYM,
	REDO_BUFFER_SIZE_SYM:                m57_REDO_BUFFER_SIZE_SYM,
	REDUNDANT_SYM:                       m57_REDUNDANT_SYM,
	REFERENCES:                          m57_REFERENCES,
	REGEXP:                              m57_REGEXP,
	RELAY:                               m57_RELAY,
	RELAYLOG_SYM:                        m57_RELAYLOG_SYM,
	RELAY_LOG_FILE_SYM:                  m57_RELAY_LOG_FILE_SYM,
	RELAY_LOG_POS_SYM:                   m57_RELAY_LOG_POS_SYM,
	RELAY_THREAD:                        m57_RELAY_THREAD,
	RELEASE_SYM:                         m57_RELEASE_SYM,
	RELOAD:                              m57_RELOAD,
	REMOVE_SYM:                          m57_REMOVE_SYM,
	RENAME:                              m57_RENAME,
	REORGANIZE_SYM:                      m57_REORGANIZE_SYM,
	REPAIR:                              m57_REPAIR,
	REPEATABLE_SYM:                      m57_REPEATABLE_SYM,
	REPEAT_SYM:                          m57_REPEAT_SYM,
	REPLACE_SYM:                         m57_REPLACE,
	REPLICATION:                         m57_REPLICATION,
	REPLICATE_DO_DB:                     m57_REPLICATE_DO_DB,
	REPLICATE_IGNORE_DB:                 m57_REPLICATE_IGNORE_DB,
	REPLICATE_DO_TABLE:                  m57_REPLICATE_DO_TABLE,
	REPLICATE_IGNORE_TABLE:              m57_REPLICATE_IGNORE_TABLE,
	REPLICATE_WILD_DO_TABLE:             m57_REPLICATE_WILD_DO_TABLE,
	REPLICATE_WILD_IGNORE_TABLE:         m57_REPLICATE_WILD_IGNORE_TABLE,
	REPLICATE_REWRITE_DB:                m57_REPLICATE_REWRITE_DB,
	REQUIRE_SYM:                         m57_REQUIRE_SYM,
	RESET_SYM:                           m57_RESET_SYM,
	RESIGNAL_SYM:                        m57_RESIGNAL_SYM,
	RESOURCES:                           m57_RESOURCES,
	RESTORE_SYM:                         m57_RESTORE_SYM,
	RESTRICT:                            m57_RESTRICT,
	RESUME_SYM:                          m57_RESUME_SYM,
	RETURNED_SQLSTATE_SYM:               m57_RETURNED_SQLSTATE_SYM,
	RETURNS_SYM:                         m57_RETURNS_SYM,
	RETURN_SYM:                          m57_RETURN_SYM,
	REVERSE_SYM:                         m57_REVERSE_SYM,
	REVOKE:                              m57_REVOKE,
	RIGHT:                               m57_RIGHT,
	ROLLBACK_SYM:                        m57_ROLLBACK_SYM,
	ROLLUP_SYM:                          m57_ROLLUP_SYM,
	ROTATE_SYM:                          m57_ROTATE_SYM,
	ROUTINE_SYM:                         m57_ROUTINE_SYM,
	ROWS_SYM:                            m57_ROWS_SYM,
	ROW_FORMAT_SYM:                      m57_ROW_FORMAT_SYM,
	ROW_SYM:                             m57_ROW_SYM,
	ROW_COUNT_SYM:                       m57_ROW_COUNT_SYM,
	RTREE_SYM:                           m57_RTREE_SYM,
	SAVEPOINT_SYM:                       m57_SAVEPOINT_SYM,
	SCHEDULE_SYM:                        m57_SCHEDULE_SYM,
	SCHEMA_NAME_SYM:                     m57_SCHEMA_NAME_SYM,
	SECOND_MICROSECOND_SYM:              m57_SECOND_MICROSECOND_SYM,
	SECOND_SYM:                          m57_SECOND_SYM,
	SECURITY_SYM:                        m57_SECURITY_SYM,
	SELECT_SYM:                          m57_SELECT_SYM,
	SENSITIVE_SYM:                       m57_SENSITIVE_SYM,
	SEPARATOR_SYM:                       m57_SEPARATOR_SYM,
	SERIALIZABLE_SYM:                    m57_SERIALIZABLE_SYM,
	SERIAL_SYM:                          m57_SERIAL_SYM,
	SESSION_SYM:                         m57_SESSION_SYM,
	SERVER_SYM:                          m57_SERVER_SYM,
	OBSOLETE_TOKEN_755:                  m57_SERVER_OPTIONS,
	SET_SYM:                             m57_SET,
	SET_VAR:                             m57_SET_VAR,
	SHARE_SYM:                           m57_SHARE_SYM,
	SHIFT_LEFT:                          m57_SHIFT_LEFT,
	SHIFT_RIGHT:                         m57_SHIFT_RIGHT,
	SHOW:                                m57_SHOW,
	SHUTDOWN:                            m57_SHUTDOWN,
	SIGNAL_SYM:                          m57_SIGNAL_SYM,
	SIGNED_SYM:                          m57_SIGNED_SYM,
	SIMPLE_SYM:                          m57_SIMPLE_SYM,
	SLAVE:                               m57_SLAVE,
	SLOW:                                m57_SLOW,
	SMALLINT_SYM:                        m57_SMALLINT,
	SNAPSHOT_SYM:                        m57_SNAPSHOT_SYM,
	SOCKET_SYM:                          m57_SOCKET_SYM,
	SONAME_SYM:                          m57_SONAME_SYM,
	SOUNDS_SYM:                          m57_SOUNDS_SYM,
	SOURCE_SYM:                          m57_SOURCE_SYM,
	SPATIAL_SYM:                         m57_SPATIAL_SYM,
	SPECIFIC_SYM:                        m57_SPECIFIC_SYM,
	SQLEXCEPTION_SYM:                    m57_SQLEXCEPTION_SYM,
	SQLSTATE_SYM:                        m57_SQLSTATE_SYM,
	SQLWARNING_SYM:                      m57_SQLWARNING_SYM,
	SQL_AFTER_GTIDS:                     m57_SQL_AFTER_GTIDS,
	SQL_AFTER_MTS_GAPS:                  m57_SQL_AFTER_MTS_GAPS,
	SQL_BEFORE_GTIDS:                    m57_SQL_BEFORE_GTIDS,
	SQL_BIG_RESULT:                      m57_SQL_BIG_RESULT,
	SQL_BUFFER_RESULT:                   m57_SQL_BUFFER_RESULT,
	OBSOLETE_TOKEN_784:                  m57_SQL_CACHE_SYM,
	SQL_CALC_FOUND_ROWS:                 m57_SQL_CALC_FOUND_ROWS,
	SQL_NO_CACHE_SYM:                    m57_SQL_NO_CACHE_SYM,
	SQL_SMALL_RESULT:                    m57_SQL_SMALL_RESULT,
	SQL_SYM:                             m57_SQL_SYM,
	SQL_THREAD:                          m57_SQL_THREAD,
	SSL_SYM:                             m57_SSL_SYM,
	STACKED_SYM:                         m57_STACKED_SYM,
	STARTING:                            m57_STARTING,
	STARTS_SYM:                          m57_STARTS_SYM,
	START_SYM:                           m57_START_SYM,
	STATS_AUTO_RECALC_SYM:               m57_STATS_AUTO_RECALC_SYM,
	STATS_PERSISTENT_SYM:                m57_STATS_PERSISTENT_SYM,
	STATS_SAMPLE_PAGES_SYM:              m57_STATS_SAMPLE_PAGES_SYM,
	STATUS_SYM:                          m57_STATUS_SYM,
	STDDEV_SAMP_SYM:                     m57_STDDEV_SAMP_SYM,
	STD_SYM:                             m57_STD_SYM,
	STOP_SYM:                            m57_STOP_SYM,
	STORAGE_SYM:                         m57_STORAGE_SYM,
	STORED_SYM:                          m57_STORED_SYM,
	STRAIGHT_JOIN:                       m57_STRAIGHT_JOIN,
	STRING_SYM:                          m57_STRING_SYM,
	SUBCLASS_ORIGIN_SYM:                 m57_SUBCLASS_ORIGIN_SYM,
	SUBDATE_SYM:                         m57_SUBDATE_SYM,
	SUBJECT_SYM:                         m57_SUBJECT_SYM,
	SUBPARTITIONS_SYM:                   m57_SUBPARTITIONS_SYM,
	SUBPARTITION_SYM:                    m57_SUBPARTITION_SYM,
	SUBSTRING:                           m57_SUBSTRING,
	SUM_SYM:                             m57_SUM_SYM,
	SUPER_SYM:                           m57_SUPER_SYM,
	SUSPEND_SYM:                         m57_SUSPEND_SYM,
	SWAPS_SYM:                           m57_SWAPS_SYM,
	SWITCHES_SYM:                        m57_SWITCHES_SYM,
	SYSDATE:                             m57_SYSDATE,
	TABLES:                              m57_TABLES,
	TABLESPACE_SYM:                      m57_TABLESPACE_SYM,
	OBSOLETE_TOKEN_820:                  m57_TABLE_REF_PRIORITY,
	TABLE_SYM:                           m57_TABLE_SYM,
	TABLE_CHECKSUM_SYM:                  m57_TABLE_CHECKSUM_SYM,
	TABLE_NAME_SYM:                      m57_TABLE_NAME_SYM,
	TEMPORARY:                           m57_TEMPORARY,
	TEMPTABLE_SYM:                       m57_TEMPTABLE_SYM,
	TERMINATED:                          m57_TERMINATED,
	TEXT_STRING:                         m57_TEXT_STRING,
	TEXT_SYM:                            m57_TEXT_SYM,
	THAN_SYM:                            m57_THAN_SYM,
	THEN_SYM:                            m57_THEN_SYM,
	TIMESTAMP_SYM:                       m57_TIMESTAMP,
	TIMESTAMP_ADD:                       m57_TIMESTAMP_ADD,
	TIMESTAMP_DIFF:                      m57_TIMESTAMP_DIFF,
	TIME_SYM:                            m57_TIME_SYM,
	TINYBLOB_SYM:                        m57_TINYBLOB,
	TINYINT_SYM:                         m57_TINYINT,
	TINYTEXT_SYN:                        m57_TINYTEXT,
	TO_SYM:                              m57_TO_SYM,
	TRAILING:                            m57_TRAILING,
	TRANSACTION_SYM:                     m57_TRANSACTION_SYM,
	TRIGGERS_SYM:                        m57_TRIGGERS_SYM,
	TRIGGER_SYM:                         m57_TRIGGER_SYM,
	TRIM:                                m57_TRIM,
	TRUE_SYM:                            m57_TRUE_SYM,
	TRUNCATE_SYM:                        m57_TRUNCATE_SYM,
	TYPES_SYM:                           m57_TYPES_SYM,
	TYPE_SYM:                            m57_TYPE_SYM,
	OBSOLETE_TOKEN_848:                  m57_UDF_RETURNS_SYM,
	ULONGLONG_NUM:                       m57_ULONGLONG_NUM,
	UNCOMMITTED_SYM:                     m57_UNCOMMITTED_SYM,
	UNDEFINED_SYM:                       m57_UNDEFINED_SYM,
	UNDERSCORE_CHARSET:                  m57_UNDERSCORE_CHARSET,
	UNDOFILE_SYM:                        m57_UNDOFILE_SYM,
	UNDO_BUFFER_SIZE_SYM:                m57_UNDO_BUFFER_SIZE_SYM,
	UNDO_SYM:                            m57_UNDO_SYM,
	UNICODE_SYM:                         m57_UNICODE_SYM,
	UNINSTALL_SYM:                       m57_UNINSTALL_SYM,
	UNION_SYM:                           m57_UNION_SYM,
	UNIQUE_SYM:                          m57_UNIQUE_SYM,
	UNKNOWN_SYM:                         m57_UNKNOWN_SYM,
	UNLOCK_SYM:                          m57_UNLOCK_SYM,
	UNSIGNED_SYM:                        m57_UNSIGNED,
	UNTIL_SYM:                           m57_UNTIL_SYM,
	UPDATE_SYM:                          m57_UPDATE_SYM,
	UPGRADE_SYM:                         m57_UPGRADE_SYM,
	USAGE:                               m57_USAGE,
	USER:                                m57_USER,
	USE_FRM:                             m57_USE_FRM,
	USE_SYM:                             m57_USE_SYM,
	USING:                               m57_USING,
	UTC_DATE_SYM:                        m57_UTC_DATE_SYM,
	UTC_TIMESTAMP_SYM:                   m57_UTC_TIMESTAMP_SYM,
	UTC_TIME_SYM:                        m57_UTC_TIME_SYM,
	VALIDATION_SYM:                      m57_VALIDATION_SYM,
	VALUES:                              m57_VALUES,
	VALUE_SYM:                           m57_VALUE_SYM,
	VARBINARY_SYM:                       m57_VARBINARY,
	VARCHAR_SYM:                         m57_VARCHAR,
	VARIABLES:                           m57_VARIABLES,
	VARIANCE_SYM:                        m57_VARIANCE_SYM,
	VARYING:                             m57_VARYING,
	VAR_SAMP_SYM:                        m57_VAR_SAMP_SYM,
	VIEW_SYM:                            m57_VIEW_SYM,
	VIRTUAL_SYM:                         m57_VIRTUAL_SYM,
	WAIT_SYM:                            m57_WAIT_SYM,
	WARNINGS:                            m57_WARNINGS,
	WEEK_SYM:                            m57_WEEK_SYM,
	WEIGHT_STRING_SYM:                   m57_WEIGHT_STRING_SYM,
	WHEN_SYM:                            m57_WHEN_SYM,
	WHERE:                               m57_WHERE,
	WHILE_SYM:                           m57_WHILE_SYM,
	WITH:                                m57_WITH,
	OBSOLETE_TOKEN_893:                  m57_WITH_CUBE_SYM,
	WITH_ROLLUP_SYM:                     m57_WITH_ROLLUP_SYM,
	WITHOUT_SYM:                         m57_WITHOUT_SYM,
	WORK_SYM:                            m57_WORK_SYM,
	WRAPPER_SYM:                         m57_WRAPPER_SYM,
	WRITE_SYM:                           m57_WRITE_SYM,
	X509_SYM:                            m57_X509_SYM,
	XA_SYM:                              m57_XA_SYM,
	XID_SYM:                             m57_XID_SYM,
	XML_SYM:                             m57_XML_SYM,
	XOR:                                 m57_XOR,
	YEAR_MONTH_SYM:                      m57_YEAR_MONTH_SYM,
	YEAR_SYM:                            m57_YEAR_SYM,
	ZEROFILL_SYM:                        m57_ZEROFILL,
	JSON_UNQUOTED_SEPARATOR_SYM:         m57_JSON_UNQUOTED_SEPARATOR_SYM,
	PERSIST_SYM:                         m57TOK_UNUSED, // New in MySQL 8.0
	ROLE_SYM:                            m57TOK_UNUSED, // New in MySQL 8.0
	ADMIN_SYM:                           m57TOK_UNUSED, // New in MySQL 8.0
	INVISIBLE_SYM:                       m57TOK_UNUSED, // New in MySQL 8.0
	VISIBLE_SYM:                         m57TOK_UNUSED, // New in MySQL 8.0
	EXCEPT_SYM:                          m57TOK_UNUSED, // New in MySQL 8.0
	COMPONENT_SYM:                       m57TOK_UNUSED, // New in MySQL 8.0
	RECURSIVE_SYM:                       m57TOK_UNUSED, // New in MySQL 8.0
	GRAMMAR_SELECTOR_EXPR:               m57TOK_UNUSED, // New in MySQL 8.0
	GRAMMAR_SELECTOR_GCOL:               m57TOK_UNUSED, // New in MySQL 8.0
	GRAMMAR_SELECTOR_PART:               m57TOK_UNUSED, // New in MySQL 8.0
	GRAMMAR_SELECTOR_CTE:                m57TOK_UNUSED, // New in MySQL 8.0
	JSON_OBJECTAGG:                      m57_JSON_OBJECTAGG,
	JSON_ARRAYAGG:                       m57_JSON_ARRAYAGG,
	OF_SYM:                              m57TOK_UNUSED, // New in MySQL 8.0
	SKIP_SYM:                            m57TOK_UNUSED, // New in MySQL 8.0
	LOCKED_SYM:                          m57TOK_UNUSED, // New in MySQL 8.0
	NOWAIT_SYM:                          m57TOK_UNUSED, // New in MySQL 8.0
	GROUPING_SYM:                        m57TOK_UNUSED, // New in MySQL 8.0
	PERSIST_ONLY_SYM:                    m57TOK_UNUSED, // New in MySQL 8.0
	HISTOGRAM_SYM:                       m57TOK_UNUSED, // New in MySQL 8.0
	BUCKETS_SYM:                         m57TOK_UNUSED, // New in MySQL 8.0
	OBSOLETE_TOKEN_930:                  m57TOK_UNUSED, // New in MySQL 8.0
	CLONE_SYM:                           m57TOK_UNUSED, // New in MySQL 8.0
	CUME_DIST_SYM:                       m57TOK_UNUSED, // New in MySQL 8.0
	DENSE_RANK_SYM:                      m57TOK_UNUSED, // New in MySQL 8.0
	EXCLUDE_SYM:                         m57TOK_UNUSED, // New in MySQL 8.0
	FIRST_VALUE_SYM:                     m57TOK_UNUSED, // New in MySQL 8.0
	FOLLOWING_SYM:                       m57TOK_UNUSED, // New in MySQL 8.0
	GROUPS_SYM:                          m57TOK_UNUSED, // New in MySQL 8.0
	LAG_SYM:                             m57TOK_UNUSED, // New in MySQL 8.0
	LAST_VALUE_SYM:                      m57TOK_UNUSED, // New in MySQL 8.0
	LEAD_SYM:                            m57TOK_UNUSED, // New in MySQL 8.0
	NTH_VALUE_SYM:                       m57TOK_UNUSED, // New in MySQL 8.0
	NTILE_SYM:                           m57TOK_UNUSED, // New in MySQL 8.0
	NULLS_SYM:                           m57TOK_UNUSED, // New in MySQL 8.0
	OTHERS_SYM:                          m57TOK_UNUSED, // New in MySQL 8.0
	OVER_SYM:                            m57TOK_UNUSED, // New in MySQL 8.0
	PERCENT_RANK_SYM:                    m57TOK_UNUSED, // New in MySQL 8.0
	PRECEDING_SYM:                       m57TOK_UNUSED, // New in MySQL 8.0
	RANK_SYM:                            m57TOK_UNUSED, // New in MySQL 8.0
	RESPECT_SYM:                         m57TOK_UNUSED, // New in MySQL 8.0
	ROW_NUMBER_SYM:                      m57TOK_UNUSED, // New in MySQL 8.0
	TIES_SYM:                            m57TOK_UNUSED, // New in MySQL 8.0
	UNBOUNDED_SYM:                       m57TOK_UNUSED, // New in MySQL 8.0
	WINDOW_SYM:                          m57TOK_UNUSED, // New in MySQL 8.0
	EMPTY_SYM:                           m57TOK_UNUSED, // New in MySQL 8.0
	JSON_TABLE_SYM:                      m57TOK_UNUSED, // New in MySQL 8.0
	NESTED_SYM:                          m57TOK_UNUSED, // New in MySQL 8.0
	ORDINALITY_SYM:                      m57TOK_UNUSED, // New in MySQL 8.0
	PATH_SYM:                            m57TOK_UNUSED, // New in MySQL 8.0
	HISTORY_SYM:                         m57TOK_UNUSED, // New in MySQL 8.0
	REUSE_SYM:                           m57TOK_UNUSED, // New in MySQL 8.0
	SRID_SYM:                            m57TOK_UNUSED, // New in MySQL 8.0
	THREAD_PRIORITY_SYM:                 m57TOK_UNUSED, // New in MySQL 8.0
	RESOURCE_SYM:                        m57TOK_UNUSED, // New in MySQL 8.0
	SYSTEM_SYM:                          m57TOK_UNUSED, // New in MySQL 8.0
	VCPU_SYM:                            m57TOK_UNUSED, // New in MySQL 8.0
	OBSOLETE_TOKEN_966:                  m57TOK_UNUSED, // New in MySQL 8.0
	OBSOLETE_TOKEN_967:                  m57TOK_UNUSED, // New in MySQL 8.0
	RESTART_SYM:                         m57TOK_UNUSED, // New in MySQL 8.0
	DEFINITION_SYM:                      m57TOK_UNUSED, // New in MySQL 8.0
	DESCRIPTION_SYM:                     m57TOK_UNUSED, // New in MySQL 8.0
	ORGANIZATION_SYM:                    m57TOK_UNUSED, // New in MySQL 8.0
	REFERENCE_SYM:                       m57TOK_UNUSED, // New in MySQL 8.0
	ACTIVE_SYM:                          m57TOK_UNUSED, // New in MySQL 8.0
	INACTIVE_SYM:                        m57TOK_UNUSED, // New in MySQL 8.0
	LATERAL_SYM:                         m57TOK_UNUSED, // New in MySQL 8.0
	ARRAY_SYM:                           m57TOK_UNUSED, // New in MySQL 8.0
	MEMBER_SYM:                          m57TOK_UNUSED, // New in MySQL 8.0
	OPTIONAL_SYM:                        m57TOK_UNUSED, // New in MySQL 8.0
	SECONDARY_SYM:                       m57TOK_UNUSED, // New in MySQL 8.0
	SECONDARY_ENGINE_SYM:                m57TOK_UNUSED, // New in MySQL 8.0
	SECONDARY_LOAD_SYM:                  m57TOK_UNUSED, // New in MySQL 8.0
	SECONDARY_UNLOAD_SYM:                m57TOK_UNUSED, // New in MySQL 8.0
	RETAIN_SYM:                          m57TOK_UNUSED, // New in MySQL 8.0
	OLD_SYM:                             m57TOK_UNUSED, // New in MySQL 8.0
	ENFORCED_SYM:                        m57TOK_UNUSED, // New in MySQL 8.0
	OJ_SYM:                              m57TOK_UNUSED, // New in MySQL 8.0
	NETWORK_NAMESPACE_SYM:               m57TOK_UNUSED, // New in MySQL 8.0
	RANDOM_SYM:                          m57TOK_UNUSED, // New in MySQL 8.0
	OBSOLETE_TOKEN_989:                  m57TOK_UNUSED, // New in MySQL 8.0
	OBSOLETE_TOKEN_990:                  m57TOK_UNUSED, // New in MySQL 8.0
	PRIVILEGE_CHECKS_USER_SYM:           m57TOK_UNUSED, // New in MySQL 8.0
	OBSOLETE_TOKEN_992:                  m57TOK_UNUSED, // New in MySQL 8.0
	REQUIRE_ROW_FORMAT_SYM:              m57TOK_UNUSED, // New in MySQL 8.0
	PASSWORD_LOCK_TIME_SYM:              m57TOK_UNUSED, // New in MySQL 8.0
	FAILED_LOGIN_ATTEMPTS_SYM:           m57TOK_UNUSED, // New in MySQL 8.0
	REQUIRE_TABLE_PRIMARY_KEY_CHECK_SYM: m57TOK_UNUSED, // New in MySQL 8.0
	STREAM_SYM:                          m57TOK_UNUSED, // New in MySQL 8.0
	OFF_SYM:                             m57TOK_UNUSED, // New in MySQL 8.0
	RETURNING_SYM:                       m57TOK_UNUSED, // New in MySQL 8.0
	MAX_EXECUTION_TIME_HINT:             m57_MAX_EXECUTION_TIME_HINT,
	RESOURCE_GROUP_HINT:                 m57TOK_UNUSED, // New in MySQL 8.0
	BKA_HINT:                            m57_BKA_HINT,
	BNL_HINT:                            m57_BNL_HINT,
	DUPSWEEDOUT_HINT:                    m57_DUPSWEEDOUT_HINT,
	FIRSTMATCH_HINT:                     m57_FIRSTMATCH_HINT,
	INTOEXISTS_HINT:                     m57_INTOEXISTS_HINT,
	LOOSESCAN_HINT:                      m57_LOOSESCAN_HINT,
	MATERIALIZATION_HINT:                m57_MATERIALIZATION_HINT,
	NO_BKA_HINT:                         m57_NO_BKA_HINT,
	NO_BNL_HINT:                         m57_NO_BNL_HINT,
	NO_ICP_HINT:                         m57_NO_ICP_HINT,
	NO_MRR_HINT:                         m57_NO_MRR_HINT,
	NO_RANGE_OPTIMIZATION_HINT:          m57_NO_RANGE_OPTIMIZATION_HINT,
	NO_SEMIJOIN_HINT:                    m57_NO_SEMIJOIN_HINT,
	MRR_HINT:                            m57_MRR_HINT,
	QB_NAME_HINT:                        m57_QB_NAME_HINT,
	SEMIJOIN_HINT:                       m57_SEMIJOIN_HINT,
	SUBQUERY_HINT:                       m57_SUBQUERY_HINT,
	DERIVED_MERGE_HINT:                  m57TOK_UNUSED, // New in MySQL 8.0
	NO_DERIVED_MERGE_HINT:               m57TOK_UNUSED, // New in MySQL 8.0
	JOIN_PREFIX_HINT:                    m57TOK_UNUSED, // New in MySQL 8.0
	JOIN_SUFFIX_HINT:                    m57TOK_UNUSED, // New in MySQL 8.0
	JOIN_ORDER_HINT:                     m57TOK_UNUSED, // New in MySQL 8.0
	JOIN_FIXED_ORDER_HINT:               m57TOK_UNUSED, // New in MySQL 8.0
	INDEX_MERGE_HINT:                    m57TOK_UNUSED, // New in MySQL 8.0
	NO_INDEX_MERGE_HINT:                 m57TOK_UNUSED, // New in MySQL 8.0
	SET_VAR_HINT:                        m57TOK_UNUSED, // New in MySQL 8.0
	SKIP_SCAN_HINT:                      m57TOK_UNUSED, // New in MySQL 8.0
	NO_SKIP_SCAN_HINT:                   m57TOK_UNUSED, // New in MySQL 8.0
	HASH_JOIN_HINT:                      m57TOK_UNUSED, // New in MySQL 8.0
	NO_HASH_JOIN_HINT:                   m57TOK_UNUSED, // New in MySQL 8.0
	HINT_ARG_NUMBER:                     m57_HINT_ARG_NUMBER,
	HINT_ARG_IDENT:                      m57_HINT_ARG_IDENT,
	HINT_ARG_QB_NAME:                    m57_HINT_ARG_QB_NAME,
	HINT_ARG_TEXT:                       m57TOK_UNUSED, // New in MySQL 8.0
	HINT_IDENT_OR_NUMBER_WITH_SCALE:     m57TOK_UNUSED, // New in MySQL 8.0
	HINT_CLOSE:                          m57_HINT_CLOSE,
	HINT_ERROR:                          m57TOK_UNUSED, // New in MySQL 8.0
	INDEX_HINT:                          m57TOK_UNUSED, // New in MySQL 8.0
	NO_INDEX_HINT:                       m57TOK_UNUSED, // New in MySQL 8.0
	JOIN_INDEX_HINT:                     m57TOK_UNUSED, // New in MySQL 8.0
	NO_JOIN_INDEX_HINT:                  m57TOK_UNUSED, // New in MySQL 8.0
	GROUP_INDEX_HINT:                    m57TOK_UNUSED, // New in MySQL 8.0
	NO_GROUP_INDEX_HINT:                 m57TOK_UNUSED, // New in MySQL 8.0
	ORDER_INDEX_HINT:                    m57TOK_UNUSED, // New in MySQL 8.0
	NO_ORDER_INDEX_HINT:                 m57TOK_UNUSED, // New in MySQL 8.0
	DERIVED_CONDITION_PUSHDOWN_HINT:     m57TOK_UNUSED, // New in MySQL 8.0
	NO_DERIVED_CONDITION_PUSHDOWN_HINT:  m57TOK_UNUSED, // New in MySQL 8.0
	HINT_ARG_FLOATING_POINT_NUMBER:      m57TOK_UNUSED, // New in MySQL 8.0
	TOK_GENERIC_VALUE:                   m57TOK_GENERIC_VALUE,
	TOK_GENERIC_VALUE_LIST:              m57TOK_GENERIC_VALUE_LIST,
	TOK_ROW_SINGLE_VALUE:                m57TOK_ROW_SINGLE_VALUE,
	TOK_ROW_SINGLE_VALUE_LIST:           m57TOK_ROW_SINGLE_VALUE_LIST,
	TOK_ROW_MULTIPLE_VALUE:              m57TOK_ROW_MULTIPLE_VALUE,
	TOK_ROW_MULTIPLE_VALUE_LIST:         m57TOK_ROW_MULTIPLE_VALUE_LIST,
	TOK_IDENT:                           m57TOK_IDENT,
	TOK_IDENT_AT:                        m57TOK_IDENT_AT,
	TOK_HINT_COMMENT_OPEN:               m57TOK_HINT_COMMENT_OPEN,
	TOK_HINT_COMMENT_CLOSE:              m57TOK_HINT_COMMENT_CLOSE,
	TOK_IN_GENERIC_VALUE_EXPRESSION:     m57TOK_UNUSED, // New in MySQL 8.0
	TOK_BY_NUMERIC_COLUMN:               m57TOK_UNUSED, // New in MySQL 8.0
	TOK_UNUSED:                          m57TOK_UNUSED, // New digest token in 8.0
	MY_SQL_PARSER_UNDEF:                 m57TOK_UNUSED, // New in MySQL 8.0
	JSON_VALUE_SYM:                      m57TOK_UNUSED, // New in MySQL 8.0
	TLS_SYM:                             m57TOK_UNUSED, // New in MySQL 8.0
	ATTRIBUTE_SYM:                       m57TOK_UNUSED, // New in MySQL 8.0
	ENGINE_ATTRIBUTE_SYM:                m57TOK_UNUSED, // New in MySQL 8.0
	SECONDARY_ENGINE_ATTRIBUTE_SYM:      m57TOK_UNUSED, // New in MySQL 8.0
	SOURCE_CONNECTION_AUTO_FAILOVER_SYM: m57TOK_UNUSED, // New in MySQL 8.0
	ZONE_SYM:                            m57TOK_UNUSED, // New in MySQL 8.0
	GRAMMAR_SELECTOR_DERIVED_EXPR:       m57TOK_UNUSED, // New in MySQL 8.0
	REPLICA_SYM:                         m57TOK_UNUSED, // New in MySQL 8.0
	REPLICAS_SYM:                        m57TOK_UNUSED, // New in MySQL 8.0
	ASSIGN_GTIDS_TO_ANONYMOUS_TRANSACTIONS_SYM: m57TOK_UNUSED, // New in MySQL 8.0
	GET_SOURCE_PUBLIC_KEY_SYM:                  m57TOK_UNUSED, // New in MySQL 8.0
	SOURCE_AUTO_POSITION_SYM:                   m57TOK_UNUSED, // New in MySQL 8.0
	SOURCE_BIND_SYM:                            m57TOK_UNUSED, // New in MySQL 8.0
	SOURCE_COMPRESSION_ALGORITHM_SYM:           m57TOK_UNUSED, // New in MySQL 8.0
	SOURCE_CONNECT_RETRY_SYM:                   m57TOK_UNUSED, // New in MySQL 8.0
	SOURCE_DELAY_SYM:                           m57TOK_UNUSED, // New in MySQL 8.0
	SOURCE_HEARTBEAT_PERIOD_SYM:                m57TOK_UNUSED, // New in MySQL 8.0
	SOURCE_HOST_SYM:                            m57TOK_UNUSED, // New in MySQL 8.0
	SOURCE_LOG_FILE_SYM:                        m57TOK_UNUSED, // New in MySQL 8.0
	SOURCE_LOG_POS_SYM:                         m57TOK_UNUSED, // New in MySQL 8.0
	SOURCE_PASSWORD_SYM:                        m57TOK_UNUSED, // New in MySQL 8.0
	SOURCE_PORT_SYM:                            m57TOK_UNUSED, // New in MySQL 8.0
	SOURCE_PUBLIC_KEY_PATH_SYM:                 m57TOK_UNUSED, // New in MySQL 8.0
	SOURCE_RETRY_COUNT_SYM:                     m57TOK_UNUSED, // New in MySQL 8.0
	SOURCE_SSL_SYM:                             m57TOK_UNUSED, // New in MySQL 8.0
	SOURCE_SSL_CA_SYM:                          m57TOK_UNUSED, // New in MySQL 8.0
	SOURCE_SSL_CAPATH_SYM:                      m57TOK_UNUSED, // New in MySQL 8.0
	SOURCE_SSL_CERT_SYM:                        m57TOK_UNUSED, // New in MySQL 8.0
	SOURCE_SSL_CIPHER_SYM:                      m57TOK_UNUSED, // New in MySQL 8.0
	SOURCE_SSL_CRL_SYM:                         m57TOK_UNUSED, // New in MySQL 8.0
	SOURCE_SSL_CRLPATH_SYM:                     m57TOK_UNUSED, // New in MySQL 8.0
	SOURCE_SSL_KEY_SYM:                         m57TOK_UNUSED, // New in MySQL 8.0
	SOURCE_SSL_VERIFY_SERVER_CERT_SYM:          m57TOK_UNUSED, // New in MySQL 8.0
	SOURCE_TLS_CIPHERSUITES_SYM:                m57TOK_UNUSED, // New in MySQL 8.0
	SOURCE_TLS_VERSION_SYM:                     m57TOK_UNUSED, // New in MySQL 8.0
	SOURCE_USER_SYM:                            m57TOK_UNUSED, // New in MySQL 8.0
	SOURCE_ZSTD_COMPRESSION_LEVEL_SYM:          m57TOK_UNUSED, // New in MySQL 8.0
	ST_COLLECT_SYM:                             m57TOK_UNUSED, // New in MySQL 8.0
	KEYRING_SYM:                                m57TOK_UNUSED, // New in MySQL 8.0
	AUTHENTICATION_SYM:                         m57TOK_UNUSED, // New in MySQL 8.0
	FACTOR_SYM:                                 m57TOK_UNUSED, // New in MySQL 8.0
	FINISH_SYM:                                 m57TOK_UNUSED, // New in MySQL 8.0
	INITIATE_SYM:                               m57TOK_UNUSED, // New in MySQL 8.0
	REGISTRATION_SYM:                           m57TOK_UNUSED, // New in MySQL 8.0
	UNREGISTER_SYM:                             m57TOK_UNUSED, // New in MySQL 8.0
	INITIAL_SYM:                                m57TOK_UNUSED, // New in MySQL 8.0
	CHALLENGE_RESPONSE_SYM:                     m57TOK_UNUSED, // New in MySQL 8.0
	GTID_ONLY_SYM:                              m57TOK_UNUSED, // New in MySQL 8.0
	INTERSECT_SYM:                              m57TOK_UNUSED, // New in MySQL 8.0
	BULK_SYM:                                   m57TOK_UNUSED, // New in MySQL 8.0
	URL_SYM:                                    m57TOK_UNUSED, // New in MySQL 8.0
	GENERATE_SYM:                               m57TOK_UNUSED, // New in MySQL 8.0
	DOLLAR_QUOTED_STRING_SYM:                   m57TOK_UNUSED, // New in MySQL 8.0
	PARSE_TREE_SYM:                             m57TOK_UNUSED, // New in MySQL 8.0
	LOG_SYM:                                    m57TOK_UNUSED, // New in MySQL 8.0
	GTIDS_SYM:                                  m57TOK_UNUSED, // New in MySQL 8.0
	PARALLEL_SYM:                               m57TOK_UNUSED, // New in MySQL 8.0
	S3_SYM:                                     m57TOK_UNUSED, // New in MySQL 8.0
	QUALIFY_SYM:                                m57TOK_UNUSED, // New in MySQL 8.0
	AUTO_SYM:                                   m57TOK_UNUSED, // New in MySQL 8.0
	MANUAL_SYM:                                 m57TOK_UNUSED, // New in MySQL 8.0
	BERNOULLI_SYM:                              m57TOK_UNUSED, // New in MySQL 8.0
	TABLESAMPLE_SYM:                            m57TOK_UNUSED, // New in MySQL 8.0
	VECTOR_SYM:                                 m57TOK_UNUSED, // New in MySQL 8.0
	PARAMETERS_SYM:                             m57TOK_UNUSED, // New in MySQL 8.0
	HEADER_SYM:                                 m57TOK_UNUSED, // New in MySQL 8.0
	LIBRARY_SYM:                                m57TOK_UNUSED, // New in MySQL 8.0
	URI_SYM:                                    m57TOK_UNUSED, // New in MySQL 8.0
	DUALITY_SYM:                                m57TOK_UNUSED, // New in MySQL 8.0
	RELATIONAL_SYM:                             m57TOK_UNUSED, // New in MySQL 8.0
	JSON_DUALITY_OBJECT_SYM:                    m57TOK_UNUSED, // New in MySQL 8.0
	ABSENT_SYM:                                 m57TOK_UNUSED, // New in MySQL 8.0
	FILE_FORMAT_SYM:                            m57TOK_UNUSED, // New in MySQL 8.0
	FILES_SYM:                                  m57TOK_UNUSED, // New in MySQL 8.0
	FILE_NAME_SYM:                              m57TOK_UNUSED, // New in MySQL 8.0
	FILE_PATTERN_SYM:                           m57TOK_UNUSED, // New in MySQL 8.0
	FILE_PREFIX_SYM:                            m57TOK_UNUSED, // New in MySQL 8.0
	ALLOW_MISSING_FILES_SYM:                    m57TOK_UNUSED, // New in MySQL 8.0
	AUTO_REFRESH_SYM:                           m57TOK_UNUSED, // New in MySQL 8.0
	AUTO_REFRESH_SOURCE_SYM:                    m57TOK_UNUSED, // New in MySQL 8.0
	VERIFY_KEY_CONSTRAINTS_SYM:                 m57TOK_UNUSED, // New in MySQL 8.0
	STRICT_LOAD_SYM:                            m57TOK_UNUSED, // New in MySQL 8.0
	EXTERNAL_FORMAT_SYM:                        m57TOK_UNUSED, // New in MySQL 8.0
	EXTERNAL_SYM:                               m57TOK_UNUSED, // New in MySQL 8.0
	MATERIALIZED_SYM:                           m57TOK_UNUSED, // New in MySQL 8.0
	GUIDED_SYM:                                 m57TOK_UNUSED, // New in MySQL 8.0
	KEYWORD_USED_AS_IDENT:                      m57TOK_UNUSED, // New in MySQL 8.0
	KEYWORD_USED_AS_KEYWORD:                    m57TOK_UNUSED, // New in MySQL 8.0
	CONDITIONLESS_JOIN:                         m57TOK_UNUSED, // New in MySQL 8.0
	PREFER_PARENTHESES:                         m57TOK_UNUSED, // New in MySQL 8.0
	EMPTY_FROM_CLAUSE:                          m57TOK_UNUSED, // New in MySQL 8.0
}
