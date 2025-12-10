# Real Image Challenge 2016

In the cinema business, a feature film is usually provided to a regional distributor based on a contract for exhibition in a particular geographical territory.

Each authorization is specified by a combination of included and excluded regions. For example, a distributor might be authorzied in the following manner:
```
Permissions for DISTRIBUTOR1
INCLUDE: INDIA
INCLUDE: UNITEDSTATES
EXCLUDE: KARNATAKA-INDIA
EXCLUDE: CHENNAI-TAMILNADU-INDIA
```
This allows `DISTRIBUTOR1` to distribute in any city inside the United States and India, *except* cities in the state of Karnataka (in India) and the city of Chennai (in Tamil Nadu, India).

At this point, asking your program if `DISTRIBUTOR1` has permission to distribute in `CHICAGO-ILLINOIS-UNITEDSTATES` should get `YES` as the answer, and asking if distribution can happen in `CHENNAI-TAMILNADU-INDIA` should of course be `NO`. Asking if distribution is possible in `BANGALORE-KARNATAKA-INDIA` should also be `NO`, because the whole state of Karnataka has been excluded.

Sometimes, a distributor might split the work of distribution amount smaller sub-distiributors inside their authorized geographies. For instance, `DISTRIBUTOR1` might assign the following permissions to `DISTRIBUTOR2`:

```
Permissions for DISTRIBUTOR2 < DISTRIBUTOR1
INCLUDE: INDIA
EXCLUDE: TAMILNADU-INDIA
```
Now, `DISTRIBUTOR2` can distribute the movie anywhere in `INDIA`, except inside `TAMILNADU-INDIA` and `KARNATAKA-INDIA` - `DISTRIBUTOR2`'s permissions are always a subset of `DISTRIBUTOR1`'s permissions. It's impossible/invalid for `DISTRIBUTOR2` to have `INCLUDE: CHINA`, for example, because `DISTRIBUTOR1` isn't authorized to do that in the first place. 

If `DISTRIBUTOR2` authorizes `DISTRIBUTOR3` to handle just the city of Hubli, Karnataka, India, for example:
```
Permissions for DISTRIBUTOR3 < DISTRIBUTOR2 < DISTRIBUTOR1
INCLUDE: HUBLI-KARNATAKA-INDIA
```
Again, `DISTRIBUTOR2` cannot authorize `DISTRIBUTOR3` with a region that they themselves do not have access to. 

## Steps to run 

```
make run
```

1. Add distributor - `./distrib add distributor1`
2. To add permissions for a distributor - `./distrib include distributor1 INDIA`
3. To check if distributor has permission - `./distrib check distributor1 ONGOLE-ANDHRAPRADESH-INDIA`
4. To list all distributors - `./distrib list`
5. To exclude a city for a distributor - `./distrib exclude distributor1 TAMILNADU-INDIA`
6. To show distributor information - `./distrib show distributor1`
7. To copy permissions from a distributor - `./distrib add distributor2 --parent distributor1`


We store the current state in `current_state.json`. If we want to reset all permissions, just remove or delete the json file. 
