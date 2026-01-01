# sports / fix the cfp

[Fix the CFP](https://fixthecfp.org) is a static site
proposing fixes to the College Football Playoff.

## The proposal

On the site, I propose a few fixes:

1. **Expand to 16 teams with no byes**: doubles on-campus playoff games,
   eliminates rust from extended layoffs
2. **Eliminate automatic qualifiers**: select the top 16 teams,
   no conference champion guarantees
3. **Replace the selection committee with the AP poll**:
   62 independent sportswriters instead of conflicted stakeholders
4. **Eliminate conference championship games**:
   crown champions by regular season record, reduce injury risk
5. **Assign bowl locations by proximity to top seeds**:
   reward top seeds with reduced travel
6. **Move recruiting window after the championship**:
   let teams focus on winning, not roster-building
7. **Move non-playoff bowls to Week 1**:
   give them a fresh identity as season openers with full rosters

## Architecture

The site uses the same pattern as my [blog](/cmd/blog).
It is a Go static site generator with HTML templates and CSS fingerprinting
hosted on Cloudflare Pages.

However, this site is a single `main.go` (~200 lines)
with one HTML template. Unlike the blog, it does not need Markdown parsing, an
articles directory, or syntax highlighting. It is just a bracket visualization
using HTML and CSS.
