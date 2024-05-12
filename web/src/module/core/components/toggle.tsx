import React from "react";

interface IToggleProps {
  label: string;
  checked?: boolean;
  setChecked: (checked: boolean) => void;
}

export const Toggle: React.FC<IToggleProps> = ({
  label,
  checked,
  setChecked,
}) => {
  return (
    <div className="flex items-center justify-center w-full">
      <label htmlFor="toogleA" className="flex items-center cursor-pointer">
        <div className="relative">
          <input
            id="toogleA"
            checked={checked}
            type="checkbox"
            className="sr-only"
            onChange={() => setChecked(!checked)}
          />
          <div className="w-14 h-8 bg-[#181E29] border-[1px] border-[#353C4A] rounded-full shadow-inner"></div>
          <div className="dot absolute w-6 h-6 bg-[#144ee3] rounded-full shadow left-1 top-1 transition"></div>
        </div>
        <div className="ml-3 text-[#C9CED6] font-light text-sm">{label}</div>
      </label>
    </div>
  );
};
