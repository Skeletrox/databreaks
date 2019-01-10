# DataBreaks
-----

DataBreaks aims to be a cross-DB, Language-independent SQL Parser that tries to decompose queries into language-independent constituents which can be parsed as per the target DB Language.


## Why

A component of this is something I am building at work, decided to make a more robust implementation that can be plugged in when/where ever needed.

## Seriously, WHY

Compiler design is interesting, and a field I have never actually explored much outside theory.


## Why Databreaks?

***Databreaks*** is a software that **Breaks** the inoperability between **Data**bases. <sup>New names and suggestions welcome.</sup>

## Can I contribute?

Please do. Fork and raise a PR, and if it aligns with my goals of this software, I shall gladly merge!

## What works?

Very, *very* little, as of now. The following table should let you know:  

| Language | Create? | Read? | Update? | Delete? | Additional Notes |
| -------- | ------- | ----- | ------- | ------- | ---------------  |
| InfluxQL | ☐ | ☑ | ☐ | ☐ | Only aggregate queries are supported, until functions can be made optional |

