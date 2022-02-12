function AppLayout({ children }) {
  return (
  <div className="min-h-screen bg-cover bg-no-repeat bg-center bg-[url('/waves1.svg')]" id="main">
    {children}
  </div>
  );
}

export default AppLayout;
