use strict;
use warnings;
use utf8;
use FindBin;
use Test::More;

my $bin = 'setlock';
my $lockfile = "$FindBin::Bin/../setlockfile-test";
BEGIN {
    `go build ./cmd/setlock/setlock.go`;
}

our @OPTS;
sub setlock {
    unshift @_, $lockfile;
    unshift @_, $_ for reverse @OPTS;
    diag "$$ setlock: @_";
    return system "./$bin @_";
}

subtest 'blocking' => sub {
    my $pid = fork;
    if (!defined $pid) {
        die "Failed to fork";
    }

    if ($pid == 0) {
        # child process
        setlock(qw(sleep 7));
        exit;
    }

    # parent process
    sleep 1; # for preparing

    my $begin = time;
    my $status = setlock(qw(echo hello));
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
        setlock(qw(sleep 3));
        exit;
    }

    # parent process
    sleep 1; # for preparing
    local @OPTS = qw(-n);
    my $status = setlock(qw(echo hello)) >> 8;
    is $status, 111, 'Execute command failed immediately';

    waitpid $pid, 0;

};

subtest 'check status code' => sub {
    {
        my $status = setlock() >> 8;
        is $status, 100;
    }
    {
        my $status = setlock("NOT_EXISTED_COMMAND_XXX") >> 8;
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
        setlock(qw(sleep 5));
        exit;
    }

    # parent process
    sleep 1; # for preparing

    my $begin = time;
    local @OPTS = qw(-n -N);
    my $status = setlock(qw(echo hello)); # -N option overwrites
    my $end = time;

    is $status, 0, 'Execute command successfully';
    ok $end - $begin > 3, 'Block rightly'; # XXX aggressive!
};

unlink $lockfile;

done_testing;

