# Origin and licence

This is a Go port of **LORD 5.02 JS**, the JavaScript rewrite of *Legend of
the Red Dragon* by Stephen Hurd (Deuce) that ships with Synchronet BBS
(`xtrn/lord` in https://gitlab.synchro.net/main/sbbs, mirrored at
https://github.com/SynchronetBBS/sbbs; the revision ported is in
`ref/synchronet/REVISION`).  Synchronet, including that directory, is
distributed under the GNU General Public License, so this port is a
derivative work under the GNU GPL as well.

`internal/assets/data` holds Synchronet's LORD text files (`lordtxt.lrd`,
looks and sayings tables, `lord.ini`) unchanged; `ref/synchronet` holds the
JavaScript that was ported, for reference.

*Legend of the Red Dragon* was created by Seth Able Robinson (Robinson
Technologies); the name and the original DOS game are owned by Metropolis
Gameport.  The original Turbo Pascal source has never been released; this
port is not derived from it.
