import { useEffect, useState } from "react";
import SelectField from "../../ui/SelectField";
import { REGIONS, STATES_BY_REGION } from "../../../constants/regions";
import {
  selectedLevel,
  selectedRegion,
  selectedState,
} from "../../../stores/dashboard";

export default function Form() {
  const [selectedLevelValue, setSelectedLevelValue] = useState("");
  const [selectedRegionValue, setSelectedRegionValue] = useState("");
  const [selectedStateValue, setSelectedStateValue] = useState("");
  const [regions, setRegions] = useState([]);
  const [states, setStates] = useState([]);

  useEffect(() => {
    if (selectedLevelValue) {
      selectedLevel.set(selectedLevelValue);
      setRegions([
        { value: "", label: "-" },
        ...REGIONS.sort((a, b) => a.label.localeCompare(b.label)),
      ]);
    } else {
      setRegions([]);
      setSelectedRegionValue("");
      setSelectedStateValue("");
    }
  }, [selectedLevelValue]);

  useEffect(() => {
    if (selectedRegionValue) {
      selectedRegion.set(selectedRegionValue);
      setStates([
        { value: "", label: "-" },
        ...Object.values(STATES_BY_REGION[selectedRegionValue])
          .flat()
          .sort((a, b) => a.label.localeCompare(b.label)),
      ]);
    } else {
      setSelectedStateValue("");
    }
  }, [selectedRegionValue]);

  useEffect(() => {
    selectedState.set(selectedStateValue);
  }, [selectedStateValue]);

  const handleChangeLevel = ({ target }) => {
    setSelectedLevelValue(target.value);
  };

  const handleChangeRegion = ({ target }) => {
    setSelectedRegionValue(target.value);
  };

  const handleChangeState = ({ target }) => {
    setSelectedStateValue(target.value);
  };

  return (
    <form className="w-full p-5 flex flex-wrap gap-5 content-center rounded-lg bg-[#121b31] border-2 border-[hsl(217,33%,20%)]">
      <SelectField
        id="level"
        label="Vizualización"
        options={[
          { value: "", label: "Nacional" },
          { value: "regional", label: "Regional" },
        ]}
        value={selectedLevelValue}
        onChange={handleChangeLevel}
      />
      {selectedLevelValue === "regional" && (
        <SelectField
          id="region"
          label="Región"
          options={regions}
          value={selectedRegionValue}
          onChange={handleChangeRegion}
        />
      )}

      {selectedRegionValue && (
        <SelectField
          id="states"
          label="Estado"
          options={states}
          value={selectedStateValue}
          onChange={handleChangeState}
        />
      )}
    </form>
  );
}
