import { Dialog, Transition } from "@headlessui/react";
import close from "../../assets/icons/close.svg";
import { Fragment } from "react";
interface propsPopup {
  children?: React.ReactNode;
  isOpen: boolean;
  onClose: () => void;
  title: string;
}

const Popup: React.FC<propsPopup> = ({ children, isOpen, onClose, title }) => {
  return (
    <Transition appear show={isOpen} as={Fragment}>
      <Dialog open={isOpen} onClose={onClose} className="relative z-50">
        <Transition.Child
          as={Fragment}
          enter="ease-out duration-[500ms]"
          enterFrom="opacity-0"
          enterTo="opacity-100"
          leave="ease-in duration-300"
          leaveFrom="opacity-100"
          leaveTo="opacity-0"
        >
          <div className="fixed inset-0 bg-[#00000663] " aria-hidden="true" />
        </Transition.Child>
        <div className="fixed inset-0  flex items-center justify-center  p-4">
          <Transition.Child
            as={Fragment}
            enter="ease-out duration-[400ms]"
            enterFrom="opacity-0 scale-95"
            enterTo="opacity-100 scale-100"
            leave="ease-in duration-300"
            leaveFrom="opacity-100 scale-100"
            leaveTo="opacity-0 scale-95"
          >
            <Dialog.Panel
              className={`flex-[0_1_49.063rem] rounded-md border-2 border-textBlockFill bg-primaryFill p-8`}
            >
              <div className=" mb-10 flex items-center justify-between ">
                <h1 className="  text-lg font-bold leading-6 tracking-[1px]">
                  {title}
                </h1>
                <button onClick={onClose}>
                  <img src={close} alt="close" />
                </button>
              </div>

              {children}
            </Dialog.Panel>
          </Transition.Child>
        </div>
      </Dialog>
    </Transition>
  );
};
export default Popup;
