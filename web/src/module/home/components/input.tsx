import React from "react";
import link from "../../../assets/link.svg";

interface IInputFormProps {
  data: string;
  setData: (data: string) => void;
  buttonType: "submit" | "button" | "reset";
}

export const InputForm: React.FC<IInputFormProps> = ({
  data,
  setData,
  buttonType,
}) => {
  return (
    <div className="flex flex-row border-[2px] border-[#C9CED6] bg-[#353C4A] rounded-3xl">
      <div className="flex justify-center items-center">
        <img src={link} alt="link" className="w-auto h-auto p-2" />
      </div>
      <div className="flex-1 overflow-hidden">
        <input
          value={data}
          onChange={(e) => setData(e.target.value)}
          className="w-full overflow-hidden h-full bg-transparent"
          placeholder="Enter the link here"
        />
      </div>
      <div className="overflow-hidden w-auto p-1 bg-transparent">
        <button
          type={buttonType}
          className="bg-[#144EE3] w-auto px-4 py-2 rounded-3xl"
        >
          Shorten now
        </button>
      </div>
    </div>
  );
};
