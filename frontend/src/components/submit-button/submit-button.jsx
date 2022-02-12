function SubmitButton({ text, onClick, d, ...rest }) {
  return (
    <button
      type="submit"
      className="relative flex items-center justify-center w-full px-4 py-3 text-xs font-semibold tracking-widest text-black uppercase transition duration-150 ease-in-out bg-gray-200 border border-gray-300 rounded-md dark:bg-gray-800 dark:border-gray-700 dark:text-white hover:bg-gray-300 dark:hover:bg-dark-gray-900 focus:outline-none focus:border-gray-800 focus:shadow-outline-gray disabled:opacity-25 group"
      onClick={onClick}
      {...rest}
    >
      <span className="absolute inset-y-0 left-0 flex items-center p-3">
        <svg
          className="w-5 h-5 text-gray-600 group-hover:text-gray-700 dark:text-gray-500 dark:group-hover:text-gray-400"
          xmlns="http://www.w3.org/2000/svg"
          viewBox="1 2 22 21"
          fill="currentColor"
          aria-hidden="true"
        >
          <path
            fillRule="evenodd"
            d={d}
            clipRule="evenodd"
          ></path>
        </svg>
      </span>
      <div className="pl-5">
        {text}
      </div>
    </button>
  );
}

export default SubmitButton;
