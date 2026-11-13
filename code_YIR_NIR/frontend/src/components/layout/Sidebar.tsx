import { NavLink } from "react-router-dom";
import styles from "./Sidebar.module.css";

interface SidebarProps {
  onNavigate?: () => void;
}

const navClass = ({ isActive }: { isActive: boolean }) =>
  isActive ? `${styles.link} ${styles.linkActive}` : styles.link;

export function Sidebar({ onNavigate }: SidebarProps) {
  return (
    <nav className={styles.nav} aria-label="Основная навигация">
      <p className={styles.sectionTitle}>Разделы</p>
      <NavLink to="/" end className={navClass} onClick={onNavigate}>
        Главная
      </NavLink>
      <NavLink to="/graph" className={navClass} onClick={onNavigate}>
        Граф
      </NavLink>
      <NavLink to="/upload" className={navClass} onClick={onNavigate}>
        Загрузка CSV
      </NavLink>
    </nav>
  );
}
