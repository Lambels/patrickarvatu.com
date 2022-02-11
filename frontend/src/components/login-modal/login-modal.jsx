import { Zoom } from "react-awesome-reveal";

function LoginModal() {
  return (
    <div className="z-50 fixed translate-x-[-50%] translate-y-[-50%] left-[50%] top-[50%]">
      <Zoom triggerOnce>
        <div className="block max-w-lg bg-white rounded-lg border-2 border-gray-200 shadow-x dark:bg-gray-800 dark:border-gray-700">
          <div className="p-5">
            <h5 className="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white">
              Login Via:
            </h5>
            <p className="font-normal text-gray-700 dark:text-gray-400 p-10">
              TODO: add custom button
            </p>
          </div>
        </div>
      </Zoom>
    </div>
  );
}

export default LoginModal;
