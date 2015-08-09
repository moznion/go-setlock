use strict;
use warnings;
use utf8;
use FindBin;
use Test::More;

my $bin;
my $lockfile = "$FindBin::Bin/../setlockfile-test";
BEGIN {
    $bin = "setlock";
    `go build -o $bin`;
}

subtest 'blocking' => sub {
    my $pid = fork;
    if (!defined $pid) {
        die "Failed to fork";
    }

    if ($pid == 0) {
        # child process
        `./$bin $lockfile sleep 7`;
        exit;
    }

    # parent process
    sleep 1; # for preparing

    my $begin = time;
    my $status = system "./$bin $lockfile echo hello";
    my $end = time;

    is $status, 0, 'Execute command successfully';
    ok $end - $begin > 4, 'Block rightly'; # XXX aggressive!
};

subtest 'non-blocking' => sub {
    my $pid = fork;
    if (!defined $pid) {
        die "Failed to fork";
    }

    if ($pid == 0) {
        # child process
        `./$bin $lockfile sleep 3`;
        exit;
    }

    # parent process
    sleep 1; # for preparing
    my $status = system("./$bin -n $lockfile echo hello") >> 8;
    is $status, 111, 'Execute command failed immediately';

    waitpid $pid, 0;

};

subtest 'check status code' => sub {
    {
        my $status = system("./$bin $lockfile") >> 8;
        is $status, 100;
    }
    {
        my $status = system("./$bin $lockfile NOT_EXISTED_COMMAND_XXX") >> 8;
        is $status, 111;
    }
};

subtest 'overwrite with negative option' => sub {
    # Blocking!

    my $pid = fork;
    if (!defined $pid) {
        die "Failed to fork";
    }

    if ($pid == 0) {
        # child process
        `./$bin $lockfile sleep 5`;
        exit;
    }

    # parent process
    sleep 1; # for preparing

    my $begin = time;
    my $status = system "./$bin -n -N $lockfile echo hello"; # -N option overwrites
    my $end = time;

    is $status, 0, 'Execute command successfully';
    ok $end - $begin > 3, 'Block rightly'; # XXX aggressive!
};

unlink $lockfile;

done_testing;

