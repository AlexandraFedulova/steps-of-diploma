import { useCallback, useEffect, useState } from "react";
import { Outlet, useLocation } from "react-router-dom";
import styles from "./AppLayout.module.css";
import { Header } from "./Header";
import { Sidebar } from "./Sidebar";

export function AppLayout() {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const location = useLocation();

  const closeSidebar = useCallback(() => setSidebarOpen(false), []);
  const toggleSidebar = useCallback(() => setSidebarOpen((open) => !open), []);

  useEffect(() => {
    closeSidebar();
  }, [location.pathname, closeSidebar]);

  useEffect(() => {
    if (!sidebarOpen) {
      return;
    }

    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        closeSidebar();
      }
    };

    window.addEventListener("keydown", onKeyDown);
    return () => window.removeEventListener("keydown", onKeyDown);
  }, [sidebarOpen, closeSidebar]);

  return (
    <div className={styles.shell}>
      <Header onMenuToggle={toggleSidebar} menuOpen={sidebarOpen} />

      <div className={styles.body}>
        <button
          type="button"
          className={`${styles.overlay} ${sidebarOpen ? styles.overlayVisible : ""}`}
          aria-label="Закрыть меню"
          onClick={closeSidebar}
          tabIndex={sidebarOpen ? 0 : -1}
        />

        <aside
          id="app-sidebar"
          className={`${styles.sidebar} ${sidebarOpen ? styles.sidebarOpen : ""}`}
        >
          <Sidebar onNavigate={closeSidebar} />
        </aside>

        <main className={styles.main}>
          <Outlet />
        </main>
      </div>
    </div>
  );
}
