import React from "react";

interface IModelProps {
  children?: React.ReactNode;
  onClose: () => void;
  title: string;
}

export const Model: React.FC<IModelProps> = ({ title, children, onClose }) => {
  return (
    <div className="fixed h-screen w-screen bg-[#9ca3af4d] top-0 left-0">
      <div className="flex justify-center items-center overflow-x-hidden overflow-y-auto fixed inset-0 z-50 outline-none focus:outline-none">
        <div className="relative my-4 mx-auto  min-w-80 md:min-w-[600px]">
          <div className="border-0 rounded-lg shadow-lg relative flex flex-col w-full bg-white outline-none focus:outline-none">
            <div className="flex items-center justify-between py-2 px-4 border-b border-solid border-gray-300 rounded-t ">
              <h3 className="font-medium text-xl">{title}</h3>
              <button
                className="bg-transparent border-0 text-black"
                onClick={onClose}
              >
                <span className="text-black opacity-7 h-8 w-8 text-xl block bg-gray-400 py-0 rounded-full">
                  x
                </span>
              </button>
            </div>
            <div className="relative p-6 flex-auto ">{children}</div>
          </div>
        </div>
      </div>
    </div>
  );
};
