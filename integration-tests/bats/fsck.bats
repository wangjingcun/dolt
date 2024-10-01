#! /usr/bin/env bats
load $BATS_TEST_DIRNAME/helper/common.bash

setup() {
    setup_no_dolt_init
}

teardown() {
    teardown_common
}

make_inserts() {
  for ((i=1; i<=25; i++))
  do
    dolt sql -q "INSERT INTO tbl (guid) VALUES (UUID())"
  done
  dolt commit -a -m "Add 25 values"
}

# Helper function to create enough chunks to allow archive to be created. Duplicate in archive.bats.
make_updates() {
  for ((i=1; i<=10; i++))
  do
    dolt sql -q	"
    SET @max_id = (SELECT MAX(i) FROM tbl);
    SET @random_id = FLOOR(1 + RAND() * @max_id);
    UPDATE tbl SET guid = UUID() WHERE i >= @random_id LIMIT 1;"
  done
  dolt commit -a -m "Update 10 values."
}

@test "fsck: bad commit" {
    mkdir ".dolt"
    cp -R "$BATS_TEST_DIRNAME/corrupt_dbs/bad_commit/" .dolt/

    dolt status

    run dolt fsck

    [ "$status" -eq 1 ]
    [[ "$output" =~ "rlmgv0komq0oj7qu4osdo759vs4c5pvg read with incorrect checksum: gpphmuvegiedtjtbfku4ru8jalfdk21u" ]]
}

# This test runs over 45 seconds, resulting in a timeout in lambdabats
# bats test_tags=no_lambda
@test "fsck: good archive" {
    dolt init
    dolt sql -q "create table tbl (i int auto_increment primary key, guid char(36))"
    dolt commit -A -m "create tbl"

    for ((j=1; j<=10; j++))
    do
        make_inserts
        make_updates
    done

    dolt gc
    dolt archive

    dolt fsck

    # NM4 - some more validation of the output would be nice.
}

@test "fsck: good journal" {
    dolt init
    dolt sql -q "create table tbl (i int auto_increment primary key, guid char(36))"
    dolt commit -Am "Create table tbl"

    make_inserts

    # Objects are in the journal. Don't gc.
    dolt fsck
}

@test "fsck: bad journal" {
    mkdir ".dolt"
    cp -R "$BATS_TEST_DIRNAME/corrupt_dbs/bad_journal/" .dolt/

    run dolt fsck
}
