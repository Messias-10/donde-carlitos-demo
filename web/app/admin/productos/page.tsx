import { ListaProductos } from "@/components/admin/ListaProductos";
import { SoloAdmin } from "@/components/admin/Sesion";

export default function AdminProductos() {
  return (
    <SoloAdmin>
      <ListaProductos />
    </SoloAdmin>
  );
}
