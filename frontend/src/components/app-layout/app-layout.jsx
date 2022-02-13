import classNames from "classnames";

function AppLayout({ children, bgConf }) {
  return (
    <div
      className={classNames({
        "min-h-screen": true,
        "bg-cover": true,
        "bg-no-repeat": true,
        "bg-center": true,
        [bgConf.bg]: true,
      })}
      id="main"
    >
      {children}
    </div>
  );
}

export default AppLayout;
