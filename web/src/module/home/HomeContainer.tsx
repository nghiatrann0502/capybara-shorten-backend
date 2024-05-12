import React, { useState } from "react";
import { TitleGradian } from "../core/components/titleGradian";
import { Toggle } from "../core/components/toggle";
import { InputForm } from "./components/input";
import { Model } from "../core/components/model";
import { useCopy } from "../core/hooks/useCopy";
import { createShotrUrl } from "./services/api";

export const HomeContainer = () => {
  const [checked, setChecked] = useState(false);
  const [url, setUrl] = useState("");
  const [model, setModel] = useState(false);
  const [result, setResult] = useState("asdf");
  const [_, copy] = useCopy();

  const handleCreateShortenLink = async (
    e: React.FormEvent<HTMLFormElement>,
  ) => {
    e.preventDefault();
    try {
      const data = await createShotrUrl(url);

      const { success, data: id } = data;
      if (success) {
        setModel(true);
        setResult(id);

        if (checked) {
          await copy(`http://localhost:5173/r/${id}`);
        }
        setUrl("");
      }
      console.log(data);
    } catch (err: any) {
      alert(err.error);
      console.log(err);
    }
  };

  return (
    <form className="mt-14 text-center" onSubmit={handleCreateShortenLink}>
      <TitleGradian
        title="Shorten Your Loooong Links :)"
        className="font-extrabold text-4xl px-5 text-center"
      />
      <div className="mt-5">
        <span className="text-sm text-[#c9ced6]">
          Linkly is an efficient and easy-to-use URL shortening service that
          streamlines your online experience.
        </span>
      </div>
      <div className="mt-5 text-center text-[#c9ced6]">
        <InputForm data={url} setData={setUrl} buttonType="submit" />
      </div>

      <div className="mt-5 flex">
        <Toggle
          label="Auto Paste from Clipboard"
          checked={checked}
          setChecked={setChecked}
        />
      </div>
      <div className="flex flex-col mt-5">
        <span className="text-[#c9ced6] text-md">
          You can create <span className="text-pink-600">05</span> more links.
        </span>
        <span className="text-[#c9ced6] text-md">
          <span className="underline text-blue-600 font-bold cursor-pointer">
            Register
          </span>{" "}
          Now to enjoy Unlimited usage
        </span>
      </div>

      {model && (
        <Model title="Your shorten url" onClose={() => setModel(false)}>
          <div className="flex flex-row justify-between border-1 bg-gray-400 border-gray-500 border-[1px] text-white px-4 py-3 rounded">
            <span>{`http://localhost:5173/r/${result}`}</span>
            <button
              type="button"
              onClick={async () => {
                const data = await copy(`http://localhost:5173/r/${result}`);
                if (data) {
                  alert("Copied");
                  setModel(false);
                }
              }}
            >
              Copy
            </button>
          </div>
        </Model>
      )}
    </form>
  );
};
