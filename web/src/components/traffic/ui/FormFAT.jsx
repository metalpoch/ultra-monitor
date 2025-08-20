import { useStore } from "@nanostores/react";
import SelectField from "../../ui/SelectField";
import {
  endDate,
  fat,
  fats,
} from "../../../stores/traffic";

export default function FormFAT() {
  const $endDate = useStore(endDate);
  const $fat = useStore(fat);
  const $fats = useStore(fats);
  const handleChangeFat = ({ target }) => {
    fat.set(target.value)
  };

  if ($fats.length > 0 && $endDate) return (
    <form className="w-full p-5 flex flex-wrap gap-5 content-center rounded-lg rounded-b-none bg-[#121b31] border-2 border-[hsl(217,33%,20%)] border-b-0">
      <SelectField
        id="fat"
        label="FAT"
        options={[
          {
            value: "",
            label: "Seleccionar FAT",
            disabled: true,
            hidden: true,
          },
          ...$fats,
        ]}
        value={$fat}
        onChange={handleChangeFat}
        disabled={$fats.length === 0}
      />
    </form>
  );
}

