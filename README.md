# GAUFRE: Go AUtomated Fitness REcorder

## How to run:

Setup your data in a file called log.json, the run the program with the following options to see your progress.

The data format expected by the json unmarshaler is shown in [this file](./data_format.json).

```bash
go build && ./gaufre "Assisted dip,Chest,Developpe militaire,Tricep pulldown,Lat pulldown"
```

The exercise names shown in the above command are not hardcoded. They are found by GAUFRE in the provided json file.

To show all available exercise names found in your json, you can use the `print_exercises.sh` in this repository to show you what exists in there and what will be a valid item for GAUFRE.

## NOTE:
Your exercise recordings nay not be in the expected format and it can be tedious to convert everything by hand.

I recommend giving your plaintext note or whatever to an AI that will do the translation to JSON for you.

You can find a prompt for an LLM that will do exactly this in `prompt.txt`.

## TODO:
- Muscle group pie chart
- toggle off total weight by default
