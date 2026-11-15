import { NuevoProducto } from "@/components/admin/FormProducto";
import { SoloAdmin } from "@/components/admin/Sesion";

export default function AdminNuevoProducto() {
  return (
    <SoloAdmin>
      <NuevoProducto />
    </SoloAdmin>
  );
}
