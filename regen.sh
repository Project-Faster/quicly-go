#!/bin/bash

echo "**************************************"
echo "*****  QUICLY-GO BINDINGS REGEN  *****"
echo "**************************************"
echo

function assert_errorcode() {
  if [[ ! "$?" -eq "0" ]]; then
    echo "********************************"
    echo "**** RESULT: FAILURE        ****"
    echo "********************************"
    exit 1
  fi
}

BASEDIR=$(dirname "$(realpath $0)")

echo [Prerequisites check: C_FOR_GO]
c-for-go -h
if [[ ! "$?" -eq "0" ]]; then
  echo [Build C-FOR-GO]
  cd deps/c-for-go
  go install -v
  assert_errorcode
fi

cd $BASEDIR
echo [Regen Errors package]
c-for-go -nostamp -nocgo -debug -path "$BASEDIR" -out quiclylib genspec/errors.linux.yml
assert_errorcode

echo [Regen Quicly Bindings package]
c-for-go -nostamp -debug -path "$BASEDIR" -out internal genspec/bindings.linux.yml
assert_errorcode

echo
echo "***************************"
echo "****  RESULT: SUCCESS  ****"
echo "***************************"
cd $BASEDIR
exit 0
