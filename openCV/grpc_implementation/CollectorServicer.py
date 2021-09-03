from grpc_basic import collect_pb2
from grpc_basic import collect_pb2_grpc
from service.takePicture import take_picture


class CollectorServicer(collect_pb2_grpc.CollectorServicer):
    """
    The server of OpenCV.
    """

    def collectImage(self, request, context):
        """
        Receive the request to collect a new image using the camera;
        store the image in S3 and return the information of that image.
        :param request: collect_pb2.Empty
        :param context: grpc._server._Context
        :return:
        """
        name, imgPath = take_picture()
        return collect_pb2.ImageInfo(name=name, path=imgPath)
